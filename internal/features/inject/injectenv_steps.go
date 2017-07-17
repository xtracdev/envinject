package inject

import (
	"os"

	"strings"

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
			Name:      aws.String("/inttest/p1"),
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
			Name:      aws.String("/inttest/p2"),
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
		os.Setenv(envinject.ParamPathEnvVar, "/inttest")
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
		os.Setenv(envinject.ParamPathEnvVar, "/inttest")
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

	Given(`^a var defined both in the param store and the environment$`, func() {
		os.Setenv("p1", "p1 from the environment")
		mixed, err = envinject.NewInjectedEnv()
		if err != nil {
			T.Errorf(err.Error())
			return
		}
	})

	When(`^I lookup the var$`, func() {
		//Lookup place in following step for convenience
	})

	Then(`^the param store value is returned$`, func() {
		p1 := mixed.Getenv("p1")
		assert.Equal(T, "p1Value", p1)
	})

	Given(`^a mixed environment$`, func() {
		//Env set above: p1 is in both, p2 is in the param store,
		//Foo and Bar are in the environment
	})

	When(`^I enumerate the vars in the environment$`, func() {
		//Done below
	})

	And(`^the same value is in both the env and the parame store$`, func() {
		//set above
	})

	Then(`^the param store vars values are returned$`, func() {
		combined := make(map[string]string)
		varsWithVals := mixed.Environ()
		for _, s := range varsWithVals {
			parts := strings.Split(s, "=")
			combined[parts[0]] = parts[1]
		}

		assert.Equal(T, "p1Value", combined["p1"])
		assert.Equal(T, "p2Value is secret", combined["p2"])
		assert.Equal(T, "Bar value", combined["Bar"])
		assert.Equal(T, "Foo value", combined["Foo"])

	})

}
