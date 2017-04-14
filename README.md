# Envinject

**NOTE** Code currently is not paging the results from describe_parameters
correctly - this needs to be corrected or only the first page will
be handled.

Optionally inject variables into the environment. Currently supports
pulling AWS SSM Parameters with a given prefix into the environment.

The goal is to be able to write code that can get its configuration from
the environment, with the values in the environment potentially injected
with values from some configuration store like the AWS SSM Parameter
store.

The InjectEnv function is the mechanism to inject AWS parameter store
values into the environment. If the AWS_PARAM_STORE_PREFIX environment 
variable is set, the code will attempt to read parameters from the 
parameter store, injecting the parameters with names statring with the 
prefix, using the name stripped of the prefix as the key. So if
AWS_PARAM_STORE_PREFIX is set to `foo-`, any variables starting with
`foo-` are injected into the environmen: `foo-var1` is injected as `var1`,
`foo-var2` is injected as `var2`, and so on. All other parameters are 
ignored.

To illustrate how this works for tasks run on an ECS cluster, a 
sample os provided. Built the sample via the provided Makefile, or just
use the image that has been pushed it docker.io as xtracdev/dumpos.

To run the sample, some set up is required. First, if the log group
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

Next, create a policy document that grants access to the parameters. Customize the 
params-access-template.json with your account number and rename it as shown in the 
command below.

<pre>
aws iam create-policy --policy-name config-params --policy-document file://param-access.json
</pre>

Next, attach the above policy to the role you created:

<pre>
aws iam attach-role-policy --role-name dumpos --policy-arn "arn:aws:iam::<account-id>:policy/config-params"
</pre>

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
time="2017-04-14T04:14:56Z" level=info msg="Looking for parameters starting with sample."
time="2017-04-14T04:14:56Z" level=info msg="Create AWS session"
time="2017-04-14T04:14:57Z" level=info msg="skipping dev.p1"
time="2017-04-14T04:14:57Z" level=info msg="skipping prod.p1"
time="2017-04-14T04:14:57Z" level=info msg="Injecting sample.1 as 1"
time="2017-04-14T04:14:57Z" level=info msg="Injecting sample.2 as 2"
time="2017-04-14T04:14:57Z" level=info msg="Injecting sample.3 as 3"
time="2017-04-14T04:14:57Z" level=info msg="Injecting sample.4 as 4"
time="2017-04-14T04:14:57Z" level=info msg="Injecting sample.5 as 5"
time="2017-04-14T04:14:57Z" level=info msg="Injecting sample.PARAM4 as PARAM4"
time="2017-04-14T04:14:57Z" level=info msg="Injecting sample.PARAM5 as PARAM5"
time="2017-04-14T04:14:57Z" level=info msg="Injecting sample.my_secret as my_secret"
time="2017-04-14T04:14:57Z" level=info msg="HOSTNAME=85336bcc9cfc"
time="2017-04-14T04:14:57Z" level=info msg="PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
time="2017-04-14T04:14:57Z" level=info msg="AWS_CONTAINER_CREDENTIALS_RELATIVE_URI=/v2/credentials/c82b2812-cef9-4b5a-8565-f8db88e5884a"
time="2017-04-14T04:14:57Z" level=info msg="AWS_PARAM_STORE_PREFIX=sample."
time="2017-04-14T04:14:57Z" level=info msg="AWS_REGION=eu-west-1"
time="2017-04-14T04:14:57Z" level=info msg="HOME=/"
time="2017-04-14T04:14:57Z" level=info msg="1=Value1"
time="2017-04-14T04:14:57Z" level=info msg="2=Value2"
time="2017-04-14T04:14:57Z" level=info msg="3=Value3"
time="2017-04-14T04:14:57Z" level=info msg="4=Value4"
time="2017-04-14T04:14:57Z" level=info msg="5=Value5"
time="2017-04-14T04:14:57Z" level=info msg="PARAM4=Param 4 Value"
time="2017-04-14T04:14:57Z" level=info msg="PARAM5=Param 5 Value"
time="2017-04-14T04:14:57Z" level=info msg="my_secret=loose lips sink ships"
</pre>

Note the decrypted read of sample.my_secret. 

## Notes

Your base container needs to have the right CA certs to allow AWS services
to be called. Consider using [scratchy](https://github.com/xtraclabs/scratchy)
for golang images, which is the scratch image plus CA certs.