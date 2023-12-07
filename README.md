# go-auth-api

## Description

Repository for sample auth api.

## Dev Settings

### Set environment variables

1. Export environment variables
    ```
    $ touch .envrc
    $ echo "dotenv" >> .envrc
    $ touch .env
    $ echo "CDK_DEFAULT_ACCOUNT=<YOUR AWS ACCOUNT FOR DEPLOY APPS>" >> .env
    $ echo "CDK_DEFAULT_REGION=<YOUR AWS REGION FOR DEPLOY APPS>" >> .env
    $ echo "AWS_VAULT_BACKEND=pass" >> .env
    $ echo "AWS_VAULT_PASS_PREFIX=aws-vault" >> .env
    $ echo "AWS_SESSION_TOKEN_TTL=1h" >> .env
    $ echo "AWS_DEFAULT_REGION=<YOUR AWS DEFAULT REGION>" >> .env
    $ direnv allow
    ```

### Install dependencies

    ```
    $ make dev-deps deploy-deps
    ```

### AWS credentials

1. Generate GPG key
    ```
    $ gpg --gen-key
    ```
1. Initialize password-store
    ```
    $ pass init <TYPED EMAIL IN GPG GENERATION>
    ```
1. Register your aws profile
    ```
    $ aws-vault add <PROFILE NAME>
    Enter Access Key ID: <YOUR AWS ACCESS KEY>
    Enter Secret Access Key: <YOUR AWS SECRET ACCESS KEY>
    ```
1. Check for successful completion
    ```
    $ aws --profile <PROFILE NAME> s3 ls
    ```

### Debug

- Visual Studio Code
    1. Get AWS credential (expires in 1h)
        ```
        $ AWS_PROFILE=<PROFILE NAME> make cache-credentials
        ```
    1. Select the Debug tab from the Activity Bar on the left side
    1. Select `auth-api-server` from the pull-down menu at the top of the Side Menu
    1. Click the green triangle icon button to start debugging

## Deploy and destroy applications

### Deploy applications with CDK in AWS

1. Build applications
    ```
    $ make build
    ```
1. Archive bin files for Lambda
    ```
    $ make archive
    ```
1. Deploy applications
    ```
    $ AWS_PROFILE=<PROFILE NAME> DEPLOY_ENV=<DEPLOY ENV NAME> make deploy
    ```

### Destroy applications with CDK in AWS

```
$ AWS_PROFILE=<PROFILE NAME> DEPLOY_ENV=<DEPLOY ENV NAME> make destroy
```

### Connect to DB via ssh tunnel

1. Fetch ssh key for bastion
    ```
    $ AWS_PROFILE=<PROFILE NAME> make fetch-bastion-key
    ```
1. Open ssh tunnel
    ```
    $ AWS_PROFILE=<PROFILE NAME> make open-bastion-tunnel
    ```
1. Connect to DB
    ```
    $ mysql -u <DB USER> -D <DB NAME> -h 127.0.0.1 -P 3307 -p
    ```
