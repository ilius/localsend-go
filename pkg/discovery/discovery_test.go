package discovery

import (
	"log/slog"
	"testing"

	"github.com/ilius/localsend-go/pkg/config"
)

func TestDiscovery(t *testing.T) {
	conf := &config.Config{}
	d := New(conf, slog.Default())
	ips, err := d.pingScan()
	if err != nil {
		t.Log(err)
	}
	t.Log(ips)
}
