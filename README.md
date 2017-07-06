# Envinject


Optionally inject variables into an environment structure that allows 
reading like a 12 factor app, but without making the parameters available
via docker inspect. This allows injecting secrets without making it super
easy to obtain the secrets from the command line.

Currently supports
pulling AWS SSM Parameters with a given prefix into the environment.

The goal is to be able to write code that can get its configuration from
the environment, with the values in the environment potentially injected
with values from some configuration store like the AWS SSM Parameter
store.

## Usage

For pass through use, simply instantiate InjectEnv with an empty
AWS_PARAM_STORE_PREFIX and use the methods on the InjectEnv type to read 
configuration set as environment variables. This is useful
when using a container workflow that does not include the use of 
AWS services.

To read parameter store variables, store the variables using a 
prefix in front of each environment variable name (which serves as
a namespace or environment tag), set AWS_PARAM_STORE_PREFIX
to the prefix, and run your app reading environment config via the
InjectEnv type methods. You'll need to configure your
environment to pick up AWS credentials, which is done in the usual
way.

Note that parameter store variables are stored in InjectEnv without
the prefix. For example, if the code needs an environment variable
P1, and the environment the code run in tags variables with the
demo prefix (and uses demo for AWS_PARAM_STORE_PREFIX), then
P1 is made available via InjectEnv Getenv("P1").

Note that only the parameter store variables with the matching prefix
are made available by the InjectEnv instance. All others are ignore.
This allows using fine grained key management to assign different 
IAM service roles to different components (if desired) to keep one
component from reading another component's secret config.

To illustrate how this works for tasks run on an ECS cluster, a 
sample is provided. Built the sample via the provided Makefile, or just
use the image that has been pushed it docker.io as xtracdev/dumpos.

To run the sample locally, use the docker command line to inject
command line arguments, e.g.

<pre>
docker run -e AWS_PARAM_STORE_PREFIX=inttest- \
 -e AWS_REGION=us-east-1 \
 -e AWS_ACCESS_KEY_ID=<access key id> \
 -e AWS_SECRET_ACCESS_KEY=<secret access key> \
 xtracdev/dumpos:latest

</pre>

Note that you can run the sample to use the local environment as well:

<pre>
docker run  \
 -e FOO=fooval \
 -e BAR=barval \
 -e BAZ=yes-its-bazval \
 xtracdev/dumpos:latest
</pre>

To run the sample on ECS, some set up is required. First, if the log group
named in the task definition template is used, create it before running
the task, e.g.

<pre>
aws logs create-log-group --log-group-name ecs-tasks
</pre>

Next, you will need to create a role and some associated policies.

First you can create the role with the baseline assume role policy
document:

<pre>
aws iam create-role --role-name dumpos --assume-role-policy-document file://ecs-tasks-trust-policy.json
</pre>

Next, create a policy document that grants access to the parameters, and permissions to
access the decryption key used to decrypt the secret parameters. Customize the 
params-access-template.json with your account number and rename it as shown in the 
command below.

<pre>
aws iam create-policy --policy-name config-params --policy-document file://param-access.json
</pre>

Next, attach the above policy to the role you created:


    aws iam attach-role-policy --role-name dumpos --policy-arn "arn:aws:iam::<account-id>:policy/config-params"


Next, seed some values to make the demo more interesting, e.g.

<pre>
aws ssm put-parameter --name sample.PARAM1 --value 'Param 1 Value' --type String
aws ssm put-parameter --name sample.PARAM2 --value 'Param 2 Value' --type String
aws ssm put-parameter --name sample.PARAM3 --value 'Param 3 Value' --type String
aws ssm put-parameter --name sample.PARAM4 --value 'Param 4 Value' --type String
aws ssm put-parameter --name sample.PARAM5 --value 'Param 5 Value' --type String
</pre>

For grins, create an encryption key, and store and encrypted parameter:

<pre>
aws kms create-key --description sample-key
aws ssm put-parameter --name sample.my_secret --value 'loose lips sink ships' --type SecureString --key-id <id of key created above>
</pre>

Copy the task definition template and customize it for your setup. Minimally
you will have to provide your own account number. When the task definition is complete,
load the definition.

<pre>
aws ecs register-task-definition --cli-input-json file://$PWD/taskdef.json
</pre>

You can the run your task

<pre>
aws ecs run-task --cluster DemoCluster --task-definition dumpos
</pre>

You should see the parameter values created above in the log output.

<pre>
time="2017-04-14T04:41:00Z" level=info msg="Looking for parameters starting with sample."
time="2017-04-14T04:41:00Z" level=info msg="Create AWS session"
time="2017-04-14T04:41:00Z" level=info msg="skipping dev.p1"
time="2017-04-14T04:41:00Z" level=info msg="skipping prod.p1"
time="2017-04-14T04:41:00Z" level=info msg="Injecting sample.1 as 1"
time="2017-04-14T04:41:00Z" level=info msg="Injecting sample.10 as 10"
time="2017-04-14T04:41:00Z" level=info msg="Injecting sample.11 as 11"
time="2017-04-14T04:41:00Z" level=info msg="Injecting sample.12 as 12"
time="2017-04-14T04:41:00Z" level=info msg="Injecting sample.13 as 13"
time="2017-04-14T04:41:00Z" level=info msg="Injecting sample.14 as 14"
time="2017-04-14T04:41:00Z" level=info msg="Injecting sample.15 as 15"
time="2017-04-14T04:41:00Z" level=info msg="Injecting sample.16 as 16"
time="2017-04-14T04:41:00Z" level=info msg="Injecting sample.17 as 17"
time="2017-04-14T04:41:00Z" level=info msg="Injecting sample.18 as 18"
time="2017-04-14T04:41:00Z" level=info msg="Injecting sample.19 as 19"
time="2017-04-14T04:41:00Z" level=info msg="Injecting sample.2 as 2"
time="2017-04-14T04:41:00Z" level=info msg="Injecting sample.20 as 20"
time="2017-04-14T04:41:00Z" level=info msg="Injecting sample.21 as 21"
time="2017-04-14T04:41:00Z" level=info msg="Injecting sample.22 as 22"
time="2017-04-14T04:41:00Z" level=info msg="Injecting sample.23 as 23"
time="2017-04-14T04:41:00Z" level=info msg="Injecting sample.24 as 24"
time="2017-04-14T04:41:00Z" level=info msg="Injecting sample.25 as 25"
time="2017-04-14T04:41:00Z" level=info msg="Injecting sample.3 as 3"
time="2017-04-14T04:41:00Z" level=info msg="Injecting sample.4 as 4"
time="2017-04-14T04:41:00Z" level=info msg="Injecting sample.5 as 5"
time="2017-04-14T04:41:00Z" level=info msg="Injecting sample.6 as 6"
time="2017-04-14T04:41:00Z" level=info msg="Injecting sample.7 as 7"
time="2017-04-14T04:41:00Z" level=info msg="Injecting sample.8 as 8"
time="2017-04-14T04:41:00Z" level=info msg="Injecting sample.9 as 9"
time="2017-04-14T04:41:00Z" level=info msg="Injecting sample.PARAM4 as PARAM4"
time="2017-04-14T04:41:00Z" level=info msg="Injecting sample.PARAM5 as PARAM5"
time="2017-04-14T04:41:00Z" level=info msg="Injecting sample.my_secret as my_secret"
time="2017-04-14T04:41:00Z" level=info msg="HOSTNAME=0dba3fee938f"
time="2017-04-14T04:41:00Z" level=info msg="PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
time="2017-04-14T04:41:00Z" level=info msg="AWS_PARAM_STORE_PREFIX=sample."
time="2017-04-14T04:41:00Z" level=info msg="AWS_REGION=eu-west-1"
time="2017-04-14T04:41:00Z" level=info msg="AWS_CONTAINER_CREDENTIALS_RELATIVE_URI=/v2/credentials/976dde59-36c8-4898-b362-e291d3eba9f3"
time="2017-04-14T04:41:00Z" level=info msg="HOME=/"
time="2017-04-14T04:41:00Z" level=info msg="1=Value1"
time="2017-04-14T04:41:00Z" level=info msg="10=Value10"
time="2017-04-14T04:41:00Z" level=info msg="11=Value11"
time="2017-04-14T04:41:00Z" level=info msg="12=Value12"
time="2017-04-14T04:41:00Z" level=info msg="13=Value13"
time="2017-04-14T04:41:00Z" level=info msg="14=Value14"
time="2017-04-14T04:41:00Z" level=info msg="15=Value15"
time="2017-04-14T04:41:00Z" level=info msg="16=Value16"
time="2017-04-14T04:41:00Z" level=info msg="17=Value17"
time="2017-04-14T04:41:00Z" level=info msg="18=Value18"
time="2017-04-14T04:41:00Z" level=info msg="19=Value19"
time="2017-04-14T04:41:00Z" level=info msg="2=Value2"
time="2017-04-14T04:41:00Z" level=info msg="20=Value20"
time="2017-04-14T04:41:00Z" level=info msg="21=Value21"
time="2017-04-14T04:41:00Z" level=info msg="22=Value22"
time="2017-04-14T04:41:00Z" level=info msg="23=Value23"
time="2017-04-14T04:41:00Z" level=info msg="24=Value24"
time="2017-04-14T04:41:00Z" level=info msg="25=Value25"
time="2017-04-14T04:41:00Z" level=info msg="3=Value3"
time="2017-04-14T04:41:00Z" level=info msg="4=Value4"
time="2017-04-14T04:41:00Z" level=info msg="5=Value5"
time="2017-04-14T04:41:00Z" level=info msg="6=Value6"
time="2017-04-14T04:41:00Z" level=info msg="7=Value7"
time="2017-04-14T04:41:00Z" level=info msg="8=Value8"
time="2017-04-14T04:41:00Z" level=info msg="9=Value9"
time="2017-04-14T04:41:00Z" level=info msg="PARAM4=Param 4 Value"
time="2017-04-14T04:41:00Z" level=info msg="PARAM5=Param 5 Value"
time="2017-04-14T04:41:00Z" level=info msg="my_secret=loose lips sink ships"
</pre>

Note the decrypted read of sample.my_secret. 

## Notes

Your base container needs to have the right CA certs to allow AWS services
to be called. Consider using [scratchy](https://github.com/xtraclabs/scratchy)
for golang images, which is the scratch image plus CA certs.


## Contributing

To contribute, you must certify you agree with the [Developer Certificate of Origin](http://developercertificate.org/)
by signing your commits via `git -s`. To create a signature, configure your user name and email address in git.
Sign with your real name, do not use pseudonyms or submit anonymous commits.


In terms of workflow:

0. For significant changes or improvement, create an issue before commencing work.
1. Fork the respository, and create a branch for your edits.
2. Add tests that cover your changes, unit tests for smaller changes, acceptance test
for more significant functionality.
3. Run gofmt on each file you change before committing your changes.
4. Run golint on each file you change before committing your changes.
5. Make sure all the tests pass before committing your changes.
6. Commit your changes and issue a pull request.

## License

(c) 2017 Fidelity Investments
Licensed under the Apache License, Version 2.0
