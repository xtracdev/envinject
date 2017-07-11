package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/xtracdev/envinject"
)

func main() {
	injected, err := envinject.NewInjectedEnv()
	if err != nil {
		log.Warn(err.Error())
	} else {
		log.Info("*** Dumping injected environment variables")
		vars := injected.Environ()
		for _, v := range vars {
			log.Info(v)
		}
	}
}
