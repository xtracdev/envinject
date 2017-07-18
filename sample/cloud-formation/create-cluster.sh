#!/bin/bash

aws cloudformation create-stack \
--stack-name pstore-ecs \
--template-url https://s3.amazonaws.com/admin-doug-deploy-us-east-1/cluster-layer.yml \
--parameters ParameterKey=BucketRoot,ParameterValue=https://s3.amazonaws.com/admin-doug-deploy-us-east-1 \
ParameterKey=KeyPairName,ParameterValue=FidoKeyPair \
ParameterKey=DeploymentColor,ParameterValue=purple \
ParameterKey=InstanceType,ParameterValue=t2.medium \
--capabilities CAPABILITY_IAM
