package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/xtraclabs/envinject"
)

func main() {
	injected, err := envinject.NewInjectedEnv()
	if err != nil {
		log.Warn(err.Error())
	} else {
		vars := injected.Environ()
		for _,v := range vars {
			log.Info(v)
		}
	}
}
