package discovery

import (
	"log/slog"

	"codeberg.org/ilius/localsend-go/pkg/config"
)

type discoveryImp struct {
	conf *config.Config
	log  *slog.Logger
}

func New(conf *config.Config, logger *slog.Logger) *discoveryImp {
	return &discoveryImp{
		conf: conf,
		log:  logger,
	}
}

// Enable broadcast and monitoring functions
func (d *discoveryImp) Start() {
	go d.listenForBroadcasts()
	go d.startBroadcast()
	go d.startHTTPBroadcast()
}
