package inject

import (
	. "github.com/gucumber/gucumber"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/aws"
	"os"
	"github.com/xtraclabs/envinject"
	"github.com/stretchr/testify/assert"
	log "github.com/Sirupsen/logrus"
)

func init() {
	var env *envinject.InjectedEnv
	Given(`^some configuration store in the SSM parameter store$`, func() {
		session, err := session.NewSession()
		if err != nil {
			T.Errorf(err.Error())
			return
		}

		svc := ssm.New(session)
		p1Input := ssm.PutParameterInput{
			Name: aws.String("inttest-p1"),
			Overwrite:aws.Bool(true),
			Type:aws.String("String"),
			Value:aws.String("p1Value"),
		}
		_, err = svc.PutParameter(&p1Input)
		if err != nil {
			T.Errorf(err.Error())
			return
		}

	})

	When(`^I create an inject environment$`, func() {
		os.Setenv("AWS_PARAM_STORE_PREFIX","inttest-")
		var err error
		env,err = envinject.NewInjectedEnv()
		if err != nil {
			T.Errorf(err.Error())
			return
		}

	})

	Then(`^I can read my environment variables based on prefix$`, func() {
		log.Info("dump env")
		all := env.Environ()
		for _,es:= range all {
			log.Info(es)
		}
		p1 := env.Getenv("p1")
		assert.Equal(T, "p1Value", p1)
	})

}
