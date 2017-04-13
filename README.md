# Envinject

Optionally inject variables into the environment. Currently supports
pulling AWS SSM Parameters with a given prefix into the environment.


aws logs create-log-group --log-group-name ecs-tasks

aws ecs register-task-definition --cli-input-json file://$PWD/taskdef.json

aws ecs run-task --cluster DemoCluster --task-definition dumpos


aws iam create-role --role-name dumpos --assume-role-policy-document file://ecs-tasks-trust-policy.json
aws iam create-policy --policy-name config-params --policy-document file://param-access.json

aws iam attach-role-policy --role-name dumpos --policy-arn "arn:aws:iam::<account-id>:policy/config-params"

aws ssm put-parameter --name sample.PARAM1 --value 'Param 1 Value' --type String
aws ssm put-parameter --name sample.PARAM2 --value 'Param 2 Value' --type String
aws ssm put-parameter --name sample.PARAM3 --value 'Param 3 Value' --type String
aws ssm put-parameter --name sample.PARAM4 --value 'Param 4 Value' --type String
aws ssm put-parameter --name sample.PARAM5 --value 'Param 5 Value' --type String