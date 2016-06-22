package memcache

import (
	"testing"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core"
	"github.com/intelsdi-x/snap/core/cdata"
	"github.com/intelsdi-x/snap/core/ctypes"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMemcachePlugin(t *testing.T) {
	Convey("Meta should return metadata for the plugin", t, func() {
		meta := Meta()
		So(meta.Name, ShouldResemble, name)
		So(meta.Version, ShouldResemble, version)
		So(meta.Type, ShouldResemble, plugin.CollectorPluginType)
	})

	Convey("Create Memcache Collector", t, func() {
		collector := NewMemcacheCollector()
		Convey("So memcache collector should not be nil", func() {
			So(collector, ShouldNotBeNil)
		})
		Convey("So memcache collector should be of Memcache type", func() {
			So(collector, ShouldHaveSameTypeAs, &Memcache{})
		})
		Convey("collector.GetConfigPolicy() should return a config policy", func() {
			configPolicy, _ := collector.GetConfigPolicy()
			Convey("So config policy should not be nil", func() {
				So(configPolicy, ShouldNotBeNil)
			})
			Convey("So config policy should be a cpolicy.ConfigPolicy", func() {
				So(configPolicy, ShouldHaveSameTypeAs, &cpolicy.ConfigPolicy{})
			})
			Convey("So config policy namespace should be /raintank/memcache", func() {
				conf := configPolicy.Get([]string{"raintank", "memcache"})
				So(conf, ShouldNotBeNil)
				So(conf.HasRules(), ShouldBeTrue)
				tables := conf.RulesAsTable()
				So(len(tables), ShouldEqual, 2)
				for _, rule := range tables {
					So(rule.Name, ShouldBeIn, "server", "proto")
					switch rule.Name {
					case "server":
						So(rule.Required, ShouldBeTrue)
						So(rule.Type, ShouldEqual, "string")
					case "proto":
						So(rule.Required, ShouldBeTrue)
						So(rule.Type, ShouldEqual, "string")
					}
				}
			})
		})
	})
}

func TestMemcacheCollectMetrics(t *testing.T) {
	cfg := setupCfg("127.0.0.1:11211", "tcp")

	Convey("Ping collector", t, func() {
		p := NewMemcacheCollector()
		mt, err := p.GetMetricTypes(cfg)
		if err != nil {
			t.Fatal("failed to get metricTypes", err)
		}
		So(len(mt), ShouldBeGreaterThan, 0)
		for _, m := range mt {
			t.Log(m.Namespace().String())
		}
		Convey("collect metrics", func() {
			mts := []plugin.MetricType{
				plugin.MetricType{
					Namespace_: core.NewNamespace(
						"raintank", "memcache", "general", "pid"),
					Config_: cfg.ConfigDataNode,
				},
			}
			metrics, err := p.CollectMetrics(mts)
			So(err, ShouldBeNil)
			So(metrics, ShouldNotBeNil)
			So(len(metrics), ShouldEqual, 1)
			So(metrics[0].Namespace()[0].Value, ShouldEqual, "raintank")
			So(metrics[0].Namespace()[1].Value, ShouldEqual, "memcache")
			for _, m := range metrics {
				So(m.Namespace()[2].Value, ShouldEqual, "general")
				So(m.Namespace()[3].Value, ShouldEqual, "pid")
			}
		})
	})
}

func setupCfg(server, proto string) plugin.ConfigType {
	node := cdata.NewNode()
	node.AddItem("server", ctypes.ConfigValueStr{Value: server})
	node.AddItem("proto", ctypes.ConfigValueStr{Value: proto})
	return plugin.ConfigType{ConfigDataNode: node}
}
