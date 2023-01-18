MODULE = $(shell go list -m)
VERSION ?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || echo "1.0.0")
PACKAGES := $(shell go list ./... | grep -v -e server -e test -e middleware -e entity -e constants -e mocks -e errors -e proto | sort -r )
LDFLAGS := -ldflags "-X main.Version=${VERSION}"

CONFIG_FILE ?= ./config/local.yml
APP_DSN ?= $(shell sed -n 's/^dsn:[[:space:]]*"\(.*\)"/\1/p' $(CONFIG_FILE))
MIGRATE := migrate -path migrations -database "$(APP_DSN)"
DOCKER_REPOSITORY := hinccvi/server
MOCKERY := mockery --name=Repository -r --output=./internal/mocks

PID_FILE := './.pid'
FSWATCH_FILE := './fswatch.cfg'

.PHONY: default
default: help

# generate help info from comments: thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help: ## help information about make commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: test
test: ## run unit tests
	@echo "mode: count" > coverage-all.out
	@$(foreach pkg,$(PACKAGES), \
		go test -p=1 -cover -covermode=count -coverprofile=coverage.out ${pkg}; \
		tail -n +2 coverage.out >> coverage-all.out;)

.PHONY: test-cover
test-cover: test ## run unit tests and show test coverage information
	go tool cover -html=coverage-all.out

.PHONY: run
run: ## run the API server
	go run ${LDFLAGS} cmd/server/main.go
	

.PHONY: run-restart
run-restart: ## restart the API server
	@pkill -P `cat $(PID_FILE)` || true
	@printf '%*s\n' "80" '' | tr ' ' -
	@echo "Source file changed. Restarting server..."
	@go run ${LDFLAGS} cmd/server/main.go & echo $$! > $(PID_FILE)
	@printf '%*s\n' "80" '' | tr ' ' -

run-live: ## run the API server with live reload support (requires fswatch)
	@go run ${LDFLAGS} cmd/server/main.go & echo $$! > $(PID_FILE)
	@fswatch -x -o --event Created --event Updated --event Renamed -r internal pkg cmd config | xargs -n1 -I {} make run-restart

.PHONY: build
build:  ## build the arm API server binary
	CGO_ENABLED=0 go build ${LDFLAGS} -a -o server $(MODULE)/cmd/server

build-amd64:  ## build the amd64 API server binary
	CGO_ENABLED=0 GOARCH=amd64 go build ${LDFLAGS} -a -o server $(MODULE)/cmd/server

.PHONY: build-docker
build-docker: ## build the program as a arm docker image
	docker build -f cmd/server/Dockerfile -t $(DOCKER_REPOSITORY):$(VERSION) .

build-docker-amd64: ## build the program as a amd64 docker image
	docker build --platform=linux/amd64 -f cmd/server/Dockerfile -t $(DOCKER_REPOSITORY):$(VERSION) .

.PHONY: push-docker
push-docker: ## push docker image to dockerhub
	docker push $(DOCKER_REPOSITORY):$(VERSION)

.PHONY: clean
clean: ## remove temporary files
	rm -rf server coverage.out coverage-all.out

.PHONY: version
version: ## display the version of the API server
	@echo $(VERSION)

.PHONY: db-start
db-start: ## start the database server
	@mkdir -p testdata/postgres
	docker run --rm --name postgres -v $(shell pwd)/testdata:/testdata \
		-v $(shell pwd)/testdata/postgres:/var/lib/postgresql/data \
		-e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=postgres -d -p 5432:5432 postgres

.PHONY: db-stop
db-stop: ## stop the database server
	docker stop postgres

.PHONY: redis-start
redis-start: ## start the redis server
	@mkdir -p $(shell pwd)/testdata/redis
	@cp  $(shell pwd)/config/redis.conf $(shell pwd)/testdata/redis/redis.conf
	docker run --rm --name redis \
		-v $(shell pwd)/testdata/redis:/usr/local/etc/redis \
		-d -p 6379:6379 redis redis-server /usr/local/etc/redis/redis.conf

.PHONY: redis-stop
redis-stop: ## stop the redis server
	docker stop redis

.PHONY: nginx-start
nginx-start: ## start the nginx server
	@mkdir -p $(shell pwd)/testdata/nginx
	@touch $(shell pwd)/testdata/nginx/access.log
	@touch $(shell pwd)/testdata/nginx/error.log
	@cp  $(shell pwd)/deployments/nginx.conf $(shell pwd)/testdata/nginx/nginx.conf
	docker run --rm --name nginx \
		-v $(shell pwd)/testdata/nginx/nginx.conf:/etc/nginx/nginx.conf:ro \
		-v $(shell pwd)/testdata/nginx:/var/log/nginx \
		-v $(shell pwd)/testdata/nginx:/var/log/nginx \
		-v $(shell pwd)/testdata/nginx:/tmp/cache \
		-dp 80:80 nginx nginx-debug -g 'daemon off;'

.PHONY: nginx-stop
nginx-stop: ## stop the nginx server
	docker stop nginx

.PHONY: testdata
testdata: ## populate the database with test data
	make migrate-reset
	@echo "Populating test data..."
	@docker exec -it postgres psql "$(APP_DSN)" -f /testdata/testdata.sql

.PHONY: lint
lint: ## run golangci-lint on all Go package (requires golangci-lint)
	@golangci-lint run

.PHONY: fmt
fmt: ## run "go fmt" on all Go packages
	@go fmt $(PACKAGES)

.PHONY: migrate
migrate: ## run all new database migrations
	@echo "Running all new database migrations..."
	@$(MIGRATE) up

.PHONY: migrate-down
migrate-down: ## revert database to the last migration step
	@echo "Reverting database to the last migration step..."
	@$(MIGRATE) down 1

.PHONY: migrate-new
migrate-new: ## create a new database migration
	@read -p "Enter the name of the new migration: " name; \
	$(MIGRATE) create -ext sql -dir migrations/ -seq $${name// /_}

.PHONY: migrate-reset
migrate-reset: ## reset database and re-run all migrations
	@echo "Resetting database..."
	@$(MIGRATE) drop
	@echo "Running all database migrations..."
	@$(MIGRATE) up

.PHONY: mockery
mockery: ## mock code autogenerator 
	@read -p "Enter the repository name: " repo; \
	struct_name="$$(tr '[:lower:]' '[:upper:]' <<< $${repo:0:1})$${repo:1}"; \
	$(MOCKERY) --dir=./internal/$${repo// /_}/repository/ --filename=$${repo// /_}Repository.go --structname=$${struct_name}Repository

.PHONY: protoc
protoc: ## compile protobuf file
	@protoc \
  -I . \
  --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  --validate_out=lang=go,paths=source_relative:. \
  -Iproto/ $$(find proto -iname "*.proto") \
  -I $$GOPATH/src/protoc-gen-validate/
