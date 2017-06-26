package envinject

import (
	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"strings"
	"fmt"
	"os"
)


const ParamPrefixEnvVar = "AWS_PARAM_STORE_PREFIX"

type InjectedEnv struct {
	passThrough bool
	environment map[string]string
}

func makePassThroughEnv() *InjectedEnv {
	injectEnv := InjectedEnv {
		passThrough:true,
		environment: make(map[string]string),
	}

	return &injectEnv
}

func makeInjectedEnv() *InjectedEnv {
	injectEnv := InjectedEnv {
		passThrough:false,
		environment: make(map[string]string),
	}

	return &injectEnv
}

func (i *InjectedEnv) InjectVar(name, value string) {
	if i.passThrough == true {
		i.environment[name] = value
	}
}


// Getenv retrieves the value of the environment variable named by the key.
// It returns the value, which will be empty if the variable is not present.
// To distinguish between an empty value and an unset value, use LookupEnv.
// Note: spec borrowed from golang.org os.Getenv
func (i *InjectedEnv) Getenv(name string) string {
	if i.passThrough == true {
		return os.Getenv(name)
	}

	return i.environment[name]
}

// LookupEnv retrieves the value of the environment variable named by the key. If
// the variable is present in the environment the value (which may be empty)
// is returned and the boolean is true. Otherwise the returned
// value will be empty and the boolean will be false.
// Note: spec borrowed from golang.org os.LookupEnv
func (i *InjectedEnv) LookupEnv(name string) (string,bool) {
	if i.passThrough == true {
		return os.LookupEnv(name)
	}

	v,ok := i.environment[name]
	return v,ok
}

func (i *InjectedEnv)  Environ() []string {
	if i.passThrough == true {
		return os.Environ()
	}

	var env []string
	for k,v := range i.environment {
		env = append(env,
			fmt.Sprintf("%s=%s", k, v),
		)
	}

	return env
}
// Environ returns a copy of strings representing the environment, in the form "key=value".
// Note: spec borrowed from golang.org os.Environ





func NewInjectedEnv() (*InjectedEnv,error) {

	//Need a parameter prefix if we are reading from the SSM parameter store
	prefix := os.Getenv(ParamPrefixEnvVar)
	if prefix == "" {
		log.Infof("%s env variable not set - reading configuration from os environment.", ParamPrefixEnvVar)
		return makePassThroughEnv(),nil
	}

	//Parameter store is indicated - create a session
	log.Infof("Looking for parameters starting with %s", prefix)

	log.Info("Create AWS session")

	sess, err := session.NewSession()
	if err != nil {
		return nil,err
	}

	//Read the params and inject them into the environment
	svc := ssm.New(sess)

	params := &ssm.DescribeParametersInput{}

	injected := makeInjectedEnv()

	for {
		resp, err := svc.DescribeParameters(params)

		if err != nil {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			return nil,err
		}

		parameterMetadata := resp.Parameters
		for _, pmd := range parameterMetadata {
			if !strings.HasPrefix(*pmd.Name, prefix) {
				log.Infof("skipping %s", *pmd.Name)
				continue
			}

			keyMinusPrefix := (*pmd.Name)[len(prefix):]
			log.Infof("Injecting %s as %s", *pmd.Name, keyMinusPrefix)

			//Retrieve parameter and inject it into the environment minus the prefix.
			params := &ssm.GetParametersInput{
				Names: []*string{
					pmd.Name,
				},
				WithDecryption: aws.Bool(true),
			}
			resp, err := svc.GetParameters(params)
			if err != nil {
				return nil,err
			}

			paramVal := resp.Parameters[0].Value
			injected.InjectVar(keyMinusPrefix, *paramVal)
		}

		nextToken := resp.NextToken
		if nextToken == nil {
			break
		}

		params = &ssm.DescribeParametersInput{
			NextToken:nextToken,
		}


	}

	return injected,nil

}