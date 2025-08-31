import boto3

# Bad practice - hardcoded credentials
aws_secret_access_key = "abcdefghijklmnopqrstuvwxyz1234567890ABCD"
AWS_SECRET_KEY = "1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZabcd"

client = boto3.client('s3')