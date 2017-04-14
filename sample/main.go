package main

import (
	log "github.com/Sirupsen/logrus"
	"os"
	"github.com/xtraclabs/envinject"
)

func main() {

	err := envinject.InjectEnv()
	if err != nil {
		log.Warn(err.Error())
	}

	vars := os.Environ()
	for _,v := range vars {
		log.Info(v)
	}
}
