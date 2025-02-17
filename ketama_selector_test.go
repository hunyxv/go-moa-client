package client

import (
	"testing"

	core "github.com/blackbeans/go-moa"
)

func TestRandom(t *testing.T) {
	nodes := []core.ServiceMeta{
		core.ServiceMeta{HostPort: "localhost:2181"},
		core.ServiceMeta{HostPort: "localhost:2182"},
		core.ServiceMeta{HostPort: "localhost:2183"}}
	strategy := NewRandomStrategy(nodes)
	host := strategy.Select("100777")
	t.Log(host)
	if host.HostPort != "localhost:2182" {
		t.Fail()
	}

	host = strategy.Select("100778")
	t.Log(host)
	if host.HostPort != "localhost:2181" {
		t.Fail()
	}

	//change
	nodes = []core.ServiceMeta{core.ServiceMeta{HostPort: "localhost:2186"}}
	strategy.ReHash(nodes)
	host = strategy.Select("100777")
	t.Log(host)
	if host.HostPort != "localhost:2186" {
		t.Fail()
	}
}

func BenchmarkRandomStrategy_Select(b *testing.B) {
	nodes := []core.ServiceMeta{
		core.ServiceMeta{HostPort: "localhost:2181"},
		core.ServiceMeta{HostPort: "localhost:2182"},
		core.ServiceMeta{HostPort: "localhost:2183"}}
	strategy := NewRandomStrategy(nodes)
	count_1 := 0
	count_2 := 0
	count_3 := 0
	for i := 0; i < b.N; i++ {
		host := strategy.Select("100777")
		if host.HostPort == "localhost:2181" {
			count_1++
		} else if host.HostPort == "localhost:2182" {
			count_2++
		} else if host.HostPort == "localhost:2183" {
			count_3++
		}
	}
	b.Logf("%d,%d,%d", count_1, count_2, count_3)
}

func TestKetamaSelector(t *testing.T) {

	nodes := []core.ServiceMeta{
		core.ServiceMeta{HostPort: "localhost:2181"},
		core.ServiceMeta{HostPort: "localhost:2182"},
		core.ServiceMeta{HostPort: "localhost:2183"}}
	strategy := NewKetamaStrategy(nodes)
	host := strategy.Select("100777")
	t.Log(host)
	if host.HostPort != "localhost:2182" {
		t.Fail()
	}

	host = strategy.Select("100778")
	t.Log(host)
	if host.HostPort != "localhost:2181" {
		t.Fail()
	}

	//change
	nodes = []core.ServiceMeta{
		core.ServiceMeta{HostPort: "localhost:2186"}}
	strategy.ReHash(nodes)
	host = strategy.Select("100777")
	t.Log(host)
	if host.HostPort != "localhost:2186" {
		t.Fail()
	}
}

func BenchmarkKetamaSelector(b *testing.B) {
	nodes := []core.ServiceMeta{
		core.ServiceMeta{HostPort: "localhost:2181"},
		core.ServiceMeta{HostPort: "localhost:2182"},
		core.ServiceMeta{HostPort: "localhost:2183"}}
	strategy := NewKetamaStrategy(nodes)
	for i := 0; i < b.N; i++ {
		host := strategy.Select("100777")
		if host.HostPort != "localhost:2182" {
			b.Fail()
		}
	}
}
