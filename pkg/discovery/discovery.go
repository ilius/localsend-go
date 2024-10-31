package discovery

import "github.com/ilius/localsend-go/pkg/config"

// Enable broadcast and monitoring functions
func Start(conf *config.Config) {
	go ListenForBroadcasts(conf)
	go StartBroadcast(conf)
	go StartHTTPBroadcast(conf)
}
