package inject

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	. "github.com/gucumber/gucumber"
	"github.com/stretchr/testify/assert"
	"github.com/xtracdev/envinject"
)

func init() {
	var env *envinject.InjectedEnv
	var mixed *envinject.InjectedEnv

	session, err := session.NewSession()
	if err != nil {
		T.Errorf(err.Error())
		return
	}

	svc := ssm.New(session)

	Given(`^some configuration store in the SSM parameter store$`, func() {
		log.Info("Assuming user has access to the default key for the account")

		p1Input := ssm.PutParameterInput{
			Name:      aws.String("inttest-p1"),
			Overwrite: aws.Bool(true),
			Type:      aws.String("String"),
			Value:     aws.String("p1Value"),
		}
		_, err = svc.PutParameter(&p1Input)
		if err != nil {
			T.Errorf(err.Error())
			return
		}

		p2Input := ssm.PutParameterInput{
			Name:      aws.String("inttest-p2"),
			Overwrite: aws.Bool(true),
			Type:      aws.String("SecureString"),
			Value:     aws.String("p2Value is secret"),
		}
		_, err = svc.PutParameter(&p2Input)
		if err != nil {
			T.Errorf(err.Error())
			return
		}

	})

	When(`^I create an inject environment$`, func() {
		os.Setenv("AWS_PARAM_STORE_PREFIX", "inttest-")
		var err error
		env, err = envinject.NewInjectedEnv()
		if err != nil {
			T.Errorf(err.Error())
			return
		}

	})

	Then(`^I can read my environment variables based on prefix$`, func() {
		log.Info("dump env")
		all := env.Environ()
		for _, es := range all {
			log.Info(es)
		}
		p1 := env.Getenv("p1")
		assert.Equal(T, "p1Value", p1)
		p2 := env.Getenv("p2")
		assert.Equal(T, "p2Value is secret", p2)
	})

	Given(`^a mix of paramstore and environment variables$`, func() {
		os.Setenv("Foo", "Foo value")
		os.Setenv("Bar", "Bar value")
	})

	When(`^I create an injected environment$`, func() {
		os.Setenv("AWS_PARAM_STORE_PREFIX", "inttest-")
		var err error
		mixed, err = envinject.NewInjectedEnv()
		if err != nil {
			T.Errorf(err.Error())
			return
		}
	})

	And(`^there are some environment vars not injected$`, func() {
		//Foo and Bar above
	})

	Then(`^I can access the non-injected variables$`, func() {
		log.Info("dump env")
		all := mixed.Environ()
		for _, es := range all {
			log.Info(es)
		}
		p1 := mixed.Getenv("p1")
		assert.Equal(T, "p1Value", p1)
		p2 := mixed.Getenv("p2")
		assert.Equal(T, "p2Value is secret", p2)
		foo := mixed.Getenv("Foo")
		assert.Equal(T, "Foo value", foo)
		bar := mixed.Getenv("Bar")
		assert.Equal(T, "Bar value", bar)
	})

}
