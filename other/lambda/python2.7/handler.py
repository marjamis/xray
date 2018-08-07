import boto3
import botocore
import os
from aws_xray_sdk.core import xray_recorder
from aws_xray_sdk.core import patch_all

patch_all()
kms_c = boto3.client('kms')
print(kms_c.list_keys())

# TODO use: http://docs.aws.amazon.com/lambda/latest/dg/python-context-object.html

def handler(event, context):
  kms_c = boto3.client('kms')
  print(kms_c.list_keys())
  sts_c = boto3.client('sts')
  sts_creds = sts_c.assume_role( RoleArn=os.getenv['role'],
    RoleSessionName='lambdaassumedrole',
    DurationSeconds=900)
  print(sts_creds)

  sns_c = boto3.client('sns', aws_access_key_id=sts_creds['Credentials']['AccessKeyId'],
    aws_secret_access_key=sts_creds['Credentials']['SecretAccessKey'],
    aws_session_token=sts_creds['Credentials']['SessionToken'])
  print(sns_c.create_topic( Name='lambdatest'))
