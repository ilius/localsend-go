package handlers

import "github.com/ilius/localsend-go/pkg/config"

var conf *config.Config

func SetConfig(confArg *config.Config) {
	conf = confArg
}
