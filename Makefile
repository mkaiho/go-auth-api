ROOT_PACKAGE:=github.com/mkaiho/go-auth-api
BIN_DIR:=_deployments/bin
SRC_DIR:=$(shell go list ./cmd/...)
BINARIES:=$(SRC_DIR:$(ROOT_PACKAGE)/%=$(BIN_DIR)/%)
ARCHIVE_DIR:=_deployments/zip
ARCHIVES:=$(SRC_DIR:$(ROOT_PACKAGE)/%=$(ARCHIVE_DIR)/%)

AWS_PROFILE ?= stage
DEPLOY_ENV ?= stage

.PHONY: build
build: clean $(BINARIES)

$(BINARIES):
	go build -o $@ $(@:$(BIN_DIR)/%=$(ROOT_PACKAGE)/%)

.PHONY: archive
archive: $(ARCHIVES)

$(ARCHIVES):$(BINARIES)
	@test -d $(ARCHIVE_DIR) || mkdir $(ARCHIVE_DIR)
	@test -d $(ARCHIVE_DIR)/cmd || mkdir $(ARCHIVE_DIR)/cmd
	@cp $(@:$(ARCHIVE_DIR)/%=$(BIN_DIR)/%) $(BIN_DIR)/bootstrap
	@zip -j $@.zip $(BIN_DIR)/bootstrap
	@rm $(BIN_DIR)/bootstrap

.PHONY: reshim
reshim:
	asdf reshim golang

.PHONY: dev-deps
dev-deps:
	go install gotest.tools/gotestsum@v1.7.0
	go install github.com/vektra/mockery/v2@latest
	@make reshim

.PHONY: deps
deps:
	go mod download

.PHONY: gen-mock
gen-mock:
	make dev-deps
	mockery --all --case underscore --recursive --keeptree

.PHONY: test
test:
	# gotestsum ./entity/... ./usecase/... ./adapter/... ./infrastructure/...
	gotestsum ./...

.PHONY: test-report
test-report:
	@rm -rf ./test-results
	@mkdir -p ./test-results
	gotestsum --junitfile ./test-results/unit-tests.xml -- -coverprofile=cover.out ./...

.PHONY: deploy-deps
deploy-deps:
	cd ./_deployments/cdk && npm i

.PHONY: cdk-test
cdk-test:
	cd ./_deployments/cdk && npm test

.PHONY: cdk-update-snapshot
cdk-update-snapshot:
	cd ./_deployments/cdk && npm test -- -u

.PHONY: deploy
deploy: cdk-test
	cd ./_deployments/cdk && npx cdk deploy --profile $(AWS_PROFILE) -c env=$(DEPLOY_ENV)

.PHONY: destroy
destroy:
	cd ./_deployments/cdk && npx cdk destroy --profile $(AWS_PROFILE) -c env=$(DEPLOY_ENV)

.PHONY: cache-credentials
cache-credentials:
	@aws-vault exec $(AWS_PROFILE) --json --prompt=terminal --duration 1h > /dev/null

.PHONY: fetch-bastion-key
fetch-bastion-key:
	eval $(shell aws --profile $(AWS_PROFILE) cloudformation describe-stacks --stack-name GoAuthApiStack | \
	jq '.Stacks[0].Outputs | select(.[].OutputKey == "goauthapigetbastionkeycommand") | "aws-vault exec $(AWS_PROFILE) -- " + .[0].OutputValue + " > bastion_key.pem"')
	chmod 0600 bastion_key.pem

.PHONY: open-bastion-tunnel
open-bastion-tunnel: fetch-bastion-key
	eval $(shell echo \
	$(shell aws --profile $(AWS_PROFILE) ec2 describe-instances \
	--filters "Name=tag:Name,Values=go-auth-api-bastion" \
	--query "Reservations[0].Instances[0].PublicDnsName" \
	--out json | \
	jq '{bastionInstanceHost:.} | @text') \
	$(shell aws --profile $(AWS_PROFILE) rds describe-db-instances \
	--filters "Name=db-instance-id,Values=go-auth-api-db" \
	--query "DBInstances[0].Endpoint" \
	--out json | \
	jq '{dbHost:.Address,dbPort:.Port} | @text') | \
	jq -s '.[0] * .[1] | "ssh -N -L 3307:"+.dbHost+":"+(.dbPort|tostring)+" -i ./bastion_key.pem -4 ec2-user@"+.bastionInstanceHost')

.PHONY: clean
clean:
	@rm -rf ${BIN_DIR}
	@rm -rf ${ARCHIVE_DIR}