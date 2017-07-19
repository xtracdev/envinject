package envinject

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"strings"
)

const ParamPathEnvVar = "AWS_PARAM_STORE_PATH"

type InjectedEnv struct {
	passThrough bool
	environment map[string]string
}

func makePassThroughEnv() *InjectedEnv {
	injectEnv := InjectedEnv{
		passThrough: true,
		environment: make(map[string]string),
	}

	return &injectEnv
}

func makeInjectedEnv() *InjectedEnv {
	injectEnv := InjectedEnv{
		passThrough: false,
		environment: make(map[string]string),
	}

	return &injectEnv
}

func (i *InjectedEnv) InjectVar(name, value string) {
	if i.passThrough != true {
		i.environment[name] = value
	}
}

// Getenv retrieves the value of the environment variable named by the key.
// It returns the value, which will be empty if the variable is not present
// in either the injected env or the external environment.
// To distinguish between an empty value and an unset value, use LookupEnv.
// Note: spec borrowed from golang.org os.Getenv
func (i *InjectedEnv) Getenv(name string) string {
	if i.passThrough == true {
		return os.Getenv(name)
	}

	v, ok := i.environment[name]
	if ok {
		return v
	}

	return os.Getenv(name)
}

// LookupEnv retrieves the value of the environment variable named by the key. If
// the variable is present in the environment the value (which may be empty)
// is returned and the boolean is true. Otherwise the returned
// value will be empty and the boolean will be false.
// Note: spec borrowed from golang.org os.LookupEnv
func (i *InjectedEnv) LookupEnv(name string) (string, bool) {
	if i.passThrough == true {
		return os.LookupEnv(name)
	}

	v, ok := i.environment[name]
	if ok {
		return v, ok
	}

	return os.LookupEnv(name)
}

// Environ returns a copy of strings representing the environment, in the form "key=value".
// Note: spec borrowed from golang.org os.Environ
func (i *InjectedEnv) Environ() []string {
	if i.passThrough == true {
		return os.Environ()
	}

	//Baseline is the environment
	env := os.Environ()

	//Overwrite with param store
	for k, v := range i.environment {
		env = append(env,
			fmt.Sprintf("%s=%s", k, v),
		)
	}

	return env
}

func NewInjectedEnv() (*InjectedEnv, error) {

	//Need a parameter path if we are reading from the SSM parameter store
	paramPath := os.Getenv(ParamPathEnvVar)
	if paramPath == "" {
		log.Infof("%s env variable not set - reading configuration from os environment.", ParamPathEnvVar)
		return makePassThroughEnv(), nil
	}

	//Parameter store is indicated - create a session
	log.Infof("Looking for parameters starting with %s", paramPath)

	log.Info("Create AWS session")

	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	//Read the params and inject them into the environment
	svc := ssm.New(sess)

	getParamsInput := &ssm.GetParametersByPathInput{
		Path:aws.String(paramPath),
		WithDecryption:aws.Bool(true),
	}

	injected := makeInjectedEnv()

	for {
		resp, err := svc.GetParametersByPath(getParamsInput)

		if err != nil {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			return nil, err
		}

		params := resp.Parameters
		for _, p := range params {
			//Guard for the prefix strip below
			if !strings.HasPrefix(*p.Name, paramPath) {
				log.Infof("skipping %s", *p.Name)
				continue
			}

			keyMinusPrefix := (*p.Name)[len(paramPath) + 1:]
			log.Infof("Injecting %s as %s", *p.Name, keyMinusPrefix)
			injected.InjectVar(keyMinusPrefix, *p.Value)
		}

		nextToken := resp.NextToken
		if nextToken == nil {
			break
		}

		getParamsInput = &ssm.GetParametersByPathInput{
			NextToken: nextToken,
		}

	}

	return injected, nil

}
