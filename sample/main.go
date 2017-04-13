package main

import (
	log "github.com/Sirupsen/logrus"
	"os"
	"github.com/xtraclabs/envinject"
	"strings"
)

func main() {
	os.Setenv(envinject.ParamPrefixEnvVar,"")
	err := envinject.InjectEnv()
	if err != nil {
		log.Warn(err.Error())
	}

	os.Setenv(envinject.ParamPrefixEnvVar,"sample")
	err = envinject.InjectEnv()
	if err != nil {
		log.Warn(err.Error())
	}

	vars := os.Environ()
	for _,v := range vars {
		if strings.HasPrefix(strings.ToUpper(v), "sample") {
			log.Info(v)
		}
	}
}
