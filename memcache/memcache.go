package memcache

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/gomemcached/client"
	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core"
	"github.com/intelsdi-x/snap/core/cdata"
	"github.com/intelsdi-x/snap/core/ctypes"
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

func (m *Memcache) CollectMetrics(mts []plugin.PluginMetricType) ([]plugin.PluginMetricType, error) {
	metrics := make([]plugin.PluginMetricType, 0)
	cfg, err := getConfig(mts[0].Config())
	if err != nil {
		return nil, err
	}

	conn, err := connect(cfg["proto"], cfg["server"])
	if err != nil {
		return nil, err
	}
	for _, m := range mts {
		ns := joinNamespace(m.Namespace())
		switch ns {
		case "raintank.memcache.general.*":
			gm, err := generalMetrics(conn)
			if err != nil {
				return metrics, err
			}
			metrics = append(metrics, gm...)
		case "raintank.memcache.settings.*":
			sm, err := settingsMetrics(conn)
			if err != nil {
				return metrics, err
			}
			metrics = append(metrics, sm...)
		case "raintank.memcache.items.*.*":
			im, err := itemsMetrics(conn)
			if err != nil {
				return metrics, err
			}
			metrics = append(metrics, im...)
		case "raintank.memcache.slabs.*.*":
			sm, err := slabsMetrics(conn)
			if err != nil {
				return metrics, err
			}
			metrics = append(metrics, sm...)
		default:
			return nil, fmt.Errorf("invalid metric requested. %s", ns)
		}
	}

	return metrics, nil
}

func (m *Memcache) GetMetricTypes(cfg plugin.PluginConfigType) ([]plugin.PluginMetricType, error) {
	mts := []plugin.PluginMetricType{}
	mts = append(mts, plugin.PluginMetricType{
		Namespace_: []string{"raintank", "memcache", "general", "*"},
		Labels_:    []core.Label{{Index: 3, Name: "stat"}},
	}, plugin.PluginMetricType{
		Namespace_: []string{"raintank", "memcache", "settings", "*"},
		Labels_:    []core.Label{{Index: 3, Name: "stat"}},
	}, plugin.PluginMetricType{
		Namespace_: []string{"raintank", "memcache", "items", "*", "*"},
		Labels_:    []core.Label{{Index: 3, Name: "slabclass"}, {Index: 4, Name: "stat"}},
	}, plugin.PluginMetricType{
		Namespace_: []string{"raintank", "memcache", "slabs", "*", "*"},
		Labels_:    []core.Label{{Index: 3, Name: "slabclass"}, {Index: 4, Name: "stat"}},
	})
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

func getConfig(cfg *cdata.ConfigDataNode) (map[string]string, error) {
	config := make(map[string]string)
	conf := cfg.Table()
	srv, ok := conf["server"]
	if !ok || srv.(ctypes.ConfigValueStr).Value == "" {
		return config, fmt.Errorf("server not defined in config.")
	}
	config["server"] = srv.(ctypes.ConfigValueStr).Value
	proto, ok := conf["proto"]
	if !ok || proto.(ctypes.ConfigValueStr).Value == "" {
		return config, fmt.Errorf("proto not defined in config.")
	}
	config["proto"] = proto.(ctypes.ConfigValueStr).Value
	return config, nil
}

func connect(proto, dest string) (*memcached.Client, error) {
	fmt.Printf("connecting to %s/%s", proto, dest)
	return memcached.Connect(proto, dest)
}

func generalMetrics(conn *memcached.Client) ([]plugin.PluginMetricType, error) {
	stats, err := conn.StatsMap("")
	metrics := make([]plugin.PluginMetricType, 0)
	if err != nil {
		return nil, err
	}
	hostname, _ := os.Hostname()
	for stat, value := range stats {
		//only include the stat if we can parse it as float64
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			continue
		}
		mt := plugin.NewPluginMetricType([]string{"raintank", "memcache", "general", stat}, time.Now(), hostname, nil, []core.Label{{Index: 3, Name: "stat"}}, v)
		metrics = append(metrics, *mt)
	}
	return metrics, nil
}

func settingsMetrics(conn *memcached.Client) ([]plugin.PluginMetricType, error) {
	stats, err := conn.StatsMap("settings")
	metrics := make([]plugin.PluginMetricType, 0)
	if err != nil {
		return nil, err
	}
	hostname, _ := os.Hostname()
	for stat, value := range stats {
		//only include the stat if we can parse it as float64
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			continue
		}
		mt := plugin.NewPluginMetricType([]string{"raintank", "memcache", "settings", stat}, time.Now(), hostname, nil, []core.Label{{Index: 3, Name: "stat"}}, v)
		metrics = append(metrics, *mt)
	}
	return metrics, nil
}

func itemsMetrics(conn *memcached.Client) ([]plugin.PluginMetricType, error) {
	stats, err := conn.StatsMap("items")
	metrics := make([]plugin.PluginMetricType, 0)
	if err != nil {
		return nil, err
	}
	hostname, _ := os.Hostname()
	for fullstat, value := range stats {
		//only include the stat if we can parse it as float64
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			continue
		}
		parts := strings.Split(fullstat, ":")
		slabclass := parts[1]
		stat := parts[2]
		mt := plugin.NewPluginMetricType([]string{"raintank", "memcache", "items", slabclass, stat}, time.Now(), hostname, nil, []core.Label{{Index: 3, Name: "slabclass"}, {Index: 4, Name: "stat"}}, v)
		metrics = append(metrics, *mt)
	}
	return metrics, nil
}

func slabsMetrics(conn *memcached.Client) ([]plugin.PluginMetricType, error) {
	stats, err := conn.StatsMap("slabs")
	metrics := make([]plugin.PluginMetricType, 0)
	if err != nil {
		return nil, err
	}
	hostname, _ := os.Hostname()
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
		mt := plugin.NewPluginMetricType([]string{"raintank", "memcache", "slabs", slabclass, stat}, time.Now(), hostname, nil, []core.Label{{Index: 3, Name: "slabclass"}, {Index: 4, Name: "stat"}}, v)
		metrics = append(metrics, *mt)
	}
	return metrics, nil
}

func joinNamespace(ns []string) string {
	return strings.Join(ns, ".")
}
