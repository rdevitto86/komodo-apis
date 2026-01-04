#!/bin/bash

# LocalStack initialization script for DynamoDB
# Creates DynamoDB tables for user and auth data

echo "Initializing DynamoDB in LocalStack..."

sleep 1

# Users table
echo "Creating Users table..."
awslocal dynamodb create-table \
  --table-name komodo-users-dev \
  --attribute-definitions \
    AttributeName=user_id,AttributeType=S \
    AttributeName=email,AttributeType=S \
  --key-schema \
    AttributeName=user_id,KeyType=HASH \
  --global-secondary-indexes \
    "IndexName=email-index,KeySchema=[{AttributeName=email,KeyType=HASH}],Projection={ProjectionType=ALL},ProvisionedThroughput={ReadCapacityUnits=5,WriteCapacityUnits=5}" \
  --provisioned-throughput \
    ReadCapacityUnits=5,WriteCapacityUnits=5 \
  2>/dev/null || echo "Users table already exists"

# User profiles table
echo "Creating UserProfiles table..."
awslocal dynamodb create-table \
  --table-name komodo-user-profiles-dev \
  --attribute-definitions \
    AttributeName=user_id,AttributeType=S \
  --key-schema \
    AttributeName=user_id,KeyType=HASH \
  --provisioned-throughput \
    ReadCapacityUnits=5,WriteCapacityUnits=5 \
  2>/dev/null || echo "UserProfiles table already exists"

# Sessions table
echo "Creating Sessions table..."
awslocal dynamodb create-table \
  --table-name komodo-sessions-dev \
  --attribute-definitions \
    AttributeName=session_id,AttributeType=S \
    AttributeName=user_id,AttributeType=S \
  --key-schema \
    AttributeName=session_id,KeyType=HASH \
  --global-secondary-indexes \
    "IndexName=user-id-index,KeySchema=[{AttributeName=user_id,KeyType=HASH}],Projection={ProjectionType=ALL},ProvisionedThroughput={ReadCapacityUnits=5,WriteCapacityUnits=5}" \
  --provisioned-throughput \
    ReadCapacityUnits=5,WriteCapacityUnits=5 \
  --stream-specification \
    StreamEnabled=true,StreamViewType=NEW_AND_OLD_IMAGES \
  2>/dev/null || echo "Sessions table already exists"

# OAuth tokens table
echo "Creating OAuthTokens table..."
awslocal dynamodb create-table \
  --table-name komodo-oauth-tokens-dev \
  --attribute-definitions \
    AttributeName=token_id,AttributeType=S \
    AttributeName=user_id,AttributeType=S \
  --key-schema \
    AttributeName=token_id,KeyType=HASH \
  --global-secondary-indexes \
    "IndexName=user-id-index,KeySchema=[{AttributeName=user_id,KeyType=HASH}],Projection={ProjectionType=ALL},ProvisionedThroughput={ReadCapacityUnits=5,WriteCapacityUnits=5}" \
  --provisioned-throughput \
    ReadCapacityUnits=5,WriteCapacityUnits=5 \
  2>/dev/null || echo "OAuthTokens table already exists"

echo "DynamoDB initialized successfully"
