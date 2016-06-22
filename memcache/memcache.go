package memcache

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/gomemcached/client"
	"github.com/intelsdi-x/snap-plugin-utilities/config"
	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core"
)

const (
	// Name of plugin
	name = "memcache"
	// Version of plugin
	version = 1
	// Type of plugin
	pluginType = plugin.CollectorPluginType
)

func Meta() *plugin.PluginMeta {
	return plugin.NewPluginMeta(name, version, pluginType, []string{plugin.SnapGOBContentType}, []string{plugin.SnapGOBContentType})
}

func NewMemcacheCollector() *Memcache {
	return &Memcache{}
}

type Memcache struct {
}

func (m *Memcache) CollectMetrics(mts []plugin.MetricType) ([]plugin.MetricType, error) {
	metrics := make([]plugin.MetricType, 0)
	cfg := getConfig(mts[0])

	conn, err := connect(cfg["proto"], cfg["server"])
	if err != nil {
		return nil, err
	}
	var gMetrics map[string]float64
	var sMetrics map[string]float64
	var iMetrics map[string]map[string]float64
	var slMetrics map[string]map[string]float64
	runTime := time.Now()
	for _, m := range mts {
		ns := m.Namespace()
		section := ns[2].Value
		switch section {
		case "general":
			stat := ns[3].Value
			if len(gMetrics) == 0 {
				gMetrics, err = generalMetrics(conn)
				if err != nil {
					return metrics, err
				}
			}

			metrics = append(metrics, plugin.MetricType{
				Data_:      gMetrics[stat],
				Namespace_: core.NewNamespace("raintank", "memcache", "general", stat),
				Timestamp_: runTime,
				Version_:   m.Version(),
			})
		case "settings":
			stat := ns[3].Value
			if len(sMetrics) == 0 {
				sMetrics, err = settingsMetrics(conn)
				if err != nil {
					return metrics, err
				}
			}
			metrics = append(metrics, plugin.MetricType{
				Data_:      sMetrics[stat],
				Namespace_: core.NewNamespace("raintank", "memcache", "settings", stat),
				Timestamp_: runTime,
				Version_:   m.Version(),
			})
		case "items":
			stat := ns[4].Value
			if len(iMetrics) == 0 {
				iMetrics, err = itemsMetrics(conn)
				if err != nil {
					return metrics, err
				}
			}
			for slab, met := range iMetrics {
				metrics = append(metrics, plugin.MetricType{
					Data_:      met[stat],
					Namespace_: core.NewNamespace("raintank", "memcache", "items", slab, stat),
					Timestamp_: runTime,
					Version_:   m.Version(),
				})
			}
		case "slabs":
			stat := ns[4].Value
			if len(slMetrics) == 0 {
				slMetrics, err = slabsMetrics(conn)
				if err != nil {
					return metrics, err
				}
			}
			for slab, met := range slMetrics {
				metrics = append(metrics, plugin.MetricType{
					Data_:      met[stat],
					Namespace_: core.NewNamespace("raintank", "memcache", "slabs", slab, stat),
					Timestamp_: runTime,
					Version_:   m.Version(),
				})
			}
		default:
			return nil, fmt.Errorf("invalid metric requested. %s", ns)
		}
	}

	return metrics, nil
}

func (m *Memcache) GetMetricTypes(pct plugin.ConfigType) ([]plugin.MetricType, error) {
	cfg := getConfig(pct)
	fmt.Printf("connecting to %s\n", cfg["server"])
	conn, err := connect(cfg["proto"], cfg["server"])
	mts := []plugin.MetricType{}
	generalm, err := generalMetrics(conn)
	if err != nil {
		return nil, err
	}
	for metricName := range generalm {
		mts = append(mts, plugin.MetricType{
			Namespace_: core.NewNamespace("raintank", "memcache", "general", metricName),
		})
	}

	settingm, err := settingsMetrics(conn)
	if err != nil {
		return nil, err
	}
	for metricName := range settingm {
		mts = append(mts, plugin.MetricType{
			Namespace_: core.NewNamespace("raintank", "memcache", "settings", metricName),
		})
	}
	itemm, err := itemsMetrics(conn)
	if err != nil {
		return nil, err
	}
	if len(itemm) > 0 {
		for _, m := range itemm {
			for metricName := range m {
				mts = append(mts, plugin.MetricType{
					Namespace_: core.NewNamespace("raintank", "memcache", "items").AddDynamicElement("slabclass", "slabclass").AddStaticElement(metricName),
				})
			}
			break
		}

	}

	slabm, err := slabsMetrics(conn)
	if err != nil {
		return nil, err
	}
	if len(slabm) > 0 {
		perSlabSeen := false
		totalsSeen := false
		for slab, m := range slabm {
			if slab == "total" {
				totalsSeen = true
				for metricName := range m {
					mts = append(mts, plugin.MetricType{
						Namespace_: core.NewNamespace("raintank", "memcache", "slabs", "total", metricName),
					})
				}
			} else {
				for metricName := range m {
					mts = append(mts, plugin.MetricType{
						Namespace_: core.NewNamespace("raintank", "memcache", "slabs").AddDynamicElement("slabclass", "slabclass").AddStaticElement(metricName),
					})
				}
				perSlabSeen = true
			}
			if perSlabSeen && totalsSeen {
				break
			}
		}
	}

	return mts, nil
}

func (m *Memcache) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	c := cpolicy.New()
	rule, _ := cpolicy.NewStringRule("proto", true)
	rule2, _ := cpolicy.NewStringRule("server", true)
	p := cpolicy.NewPolicyNode()
	p.Add(rule)
	p.Add(rule2)
	c.Add([]string{"raintank", "memcache"}, p)
	return c, nil
}

func getConfig(cfg interface{}) map[string]string {
	conf := make(map[string]string)
	items, err := config.GetConfigItems(cfg, "server", "proto")
	if err != nil {
		log.Fatal(err.Error())
	}
	conf["server"] = items["server"].(string)
	conf["proto"] = items["proto"].(string)
	return conf
}

func connect(proto, dest string) (*memcached.Client, error) {
	return memcached.Connect(proto, dest)
}

func generalMetrics(conn *memcached.Client) (map[string]float64, error) {
	stats, err := conn.StatsMap("")
	metrics := make(map[string]float64)
	if err != nil {
		return nil, err
	}
	for stat, value := range stats {
		//only include the stat if we can parse it as float64
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			continue
		}
		metrics[stat] = v
	}
	return metrics, nil
}

func settingsMetrics(conn *memcached.Client) (map[string]float64, error) {
	stats, err := conn.StatsMap("settings")
	metrics := make(map[string]float64)
	if err != nil {
		return nil, err
	}

	for stat, value := range stats {
		//only include the stat if we can parse it as float64
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			continue
		}
		metrics[stat] = v
	}
	return metrics, nil
}

func itemsMetrics(conn *memcached.Client) (map[string]map[string]float64, error) {
	stats, err := conn.StatsMap("items")
	metrics := make(map[string]map[string]float64)
	if err != nil {
		return nil, err
	}

	for fullstat, value := range stats {
		//only include the stat if we can parse it as float64
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			continue
		}
		parts := strings.Split(fullstat, ":")
		slabclass := parts[1]
		stat := parts[2]
		if _, ok := metrics[slabclass]; !ok {
			metrics[slabclass] = make(map[string]float64)
		}
		metrics[slabclass][stat] = v
	}
	return metrics, nil
}

func slabsMetrics(conn *memcached.Client) (map[string]map[string]float64, error) {
	stats, err := conn.StatsMap("slabs")
	metrics := make(map[string]map[string]float64)
	if err != nil {
		return nil, err
	}
	for fullstat, value := range stats {
		//only include the stat if we can parse it as float64
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			continue
		}
		parts := strings.Split(fullstat, ":")
		var slabclass string
		var stat string
		if len(parts) == 1 {
			slabclass = "total"
			stat = parts[0]
		} else {
			slabclass = parts[0]
			stat = parts[1]
		}
		if _, ok := metrics[slabclass]; !ok {
			metrics[slabclass] = make(map[string]float64)
		}
		metrics[slabclass][stat] = v
	}
	return metrics, nil
}
