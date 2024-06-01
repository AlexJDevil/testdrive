BUILD_TO_DIR := .bin
GO_LINUX := GOOS=linux GOARCH=amd64 CGO_ENABLED=0
export AWS_ACCOUNT=808475159191
export AWS_REGION=us-west-1

all: invoke

clean:
	rm -rf .aws-sam

lint:
	golangci-lint run	

develop-clean:
	rm -rf $(BUILD_TO_DIR)
	mkdir -p $(BUILD_TO_DIR)

develop: develop-clean
	go fmt ./...
	$(GO_LINUX) go build -o $(BUILD_TO_DIR)/bootstrap ./lambda/main.go;

invoke: develop
	sam local start-api --env-vars env.json --template testdrive-api.yaml --region us-west-1 --port 9070 --docker-network host --invoke-image amazon/aws-sam-cli-emulation-image-go1.x --skip-pull-image

#run output of this command so environment variables are set.
update-creds:	
	export $(shell printf "AWS_ACCESS_KEY_ID=%s AWS_SECRET_ACCESS_KEY=%s AWS_SESSION_TOKEN=%s" \
	$(shell aws sts assume-role \
	--profile 808475159191_AdministratorAccess \
	--role-arn arn:aws:iam::808475159191:role/AdminRole \
	--role-session-name AWSCLI-Session \
	--query "Credentials.[AccessKeyId,SecretAccessKey,SessionToken]" \
	--output text))

deploy: 
	cdk deploy
