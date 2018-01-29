#!/bin/bash

bucket=bushwack

# Package using input template
aws cloudformation package \
   --template-file template.yml \
   --output-template-file cfn-output.yaml \
   --s3-bucket $bucket

# Upload latest output to public S3 bucket
aws s3 cp cfn-output.yaml s3://$bucket
