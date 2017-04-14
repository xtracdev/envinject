package envinject

import (
	log "github.com/Sirupsen/logrus"
	"os"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"strings"
)


const ParamPrefixEnvVar = "AWS_PARAM_STORE_PREFIX"

func InjectEnv() error {

	//Need a parameter prefix if we are reading from the SSM parameter store
	prefix := os.Getenv(ParamPrefixEnvVar)
	if prefix == "" {
		log.Infof("%s env variable not set - reading configuration from os environment.", ParamPrefixEnvVar)
		return nil
	}

	//Parameter store is indicated - create a session
	log.Infof("Looking for parameters starting with %s", prefix)

	log.Info("Create AWS session")

	sess, err := session.NewSession()
	if err != nil {
		return err
	}

	//Read the params and inject them into the environment
	svc := ssm.New(sess)

	params := &ssm.DescribeParametersInput{}

	for {
		resp, err := svc.DescribeParameters(params)

		if err != nil {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			return err
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
				return err
			}

			paramVal := resp.Parameters[0].Value
			os.Setenv(keyMinusPrefix, *paramVal)
		}

		nextToken := resp.NextToken
		if nextToken == nil {
			break
		}

		params = &ssm.DescribeParametersInput{
			NextToken:nextToken,
		}


	}

	return nil
}