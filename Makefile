.EXPORT_ALL_VARIABLES:

COMPOSE_PROJECT_NAME=mm-printing-backend
DOCKER_FILE=./deployments/godev.Dockerfile
REDIS_PASSWORD=123

all: build-docker-compose run-s3-storage run-db run-swagger-ui  run-redis run-init-bucket run-migrate run-nats_streaming
build-docker-compose:
	docker-compose --file ./deployments/docker-compose.yml build
run-gateway:
	docker-compose --file ./deployments/docker-compose.yml up -d gateway
run-nats_streaming:
	docker-compose --file ./deployments/docker-compose.yml up -d nats_streaming
run-redis:
	docker-compose --file ./deployments/docker-compose.yml up -d redis
run-db:
	docker-compose --file ./deployments/docker-compose.yml up -d db
run-clean-db:
	docker-compose --file ./deployments/docker-compose.yml exec db ./cockroach sql --insecure --execute="DROP DATABASE IF EXISTS postgres CASCADE"
	docker-compose --file ./deployments/docker-compose.yml exec db ./cockroach sql --insecure --execute="CREATE DATABASE IF NOT EXISTS postgres"
run-migrate:
	docker-compose --file ./deployments/docker-compose.yml up migrate
run-hydra:
	docker-compose --file ./deployments/docker-compose.yml up hydra
run-gezu:
	docker-compose --file ./deployments/docker-compose.yml up gezu
run-aurora:
	docker-compose --file ./deployments/docker-compose.yml up aurora
run-swagger-ui:
	docker-compose --file ./deployments/docker-compose.yml up -d swagger-ui
run-jaeger:
	docker-compose --file ./deployments/docker-compose.yml up -d jaeger
run-pyroscope:
	docker-compose --file ./deployments/docker-compose.yml up -d pyroscope
run-s3-storage:
	docker-compose --file ./deployments/docker-compose.yml up -d s3_storage
run-init-bucket:
	sleep 1
	docker-compose --file ./deployments/docker-compose.yml up mc_client
test-unit:
	go test ./pkg/... -cover -covermode=count -coverprofile=cover.out -coverpkg=./pkg/...
	go tool cover -func=cover.out
db:
	 make run-clean-db && make run-migrate && go run tools/gen-model.go

TEST_FILE=.
test:
	# build dummy binnary server which can get coverage
	#go test --coverpkg="./internal/..." -coverprofile=coverage.out -c cmd/server/main.go cmd/server/main_test.go -o server_test
	# run gezu server
	#./server_test -test.coverprofile=coverage.out -test.v -test.run=TestMain gezu --configPath ./resources/configs/local/gezu.config.yaml
	# run test from external
	docker-compose --file ./deployments/docker-compose.yml exec hydra sh -c "cd features/ && \
		go test -v \
			--godog.format=pretty \
			--godog.random \
			-- $(TEST_FILE)"
	# kill server_test to get code coverage report
	#pkill server_test
	# get coverage
	#go tool cover -func=coverage.out
gen-err-code:
	go run tools/gen-err-code.go
gen-model-from-db:
	go run tools/gen-model.go
gen-mock:
	docker-compose --file ./deployments/docker-compose.yml run mock_gen $(args)
gen-mock-repo:
		echo "generating mock repo for"; \
		docker-compose --file ./deployments/docker-compose.yml run mock_gen \
		--all \
		--dir=./pkg/repository/ \
		--outpkg=mockreposiotry \
		--output=./mocks/repository \
		--keeptree; \
		echo "generating mock service for"; \
		docker-compose --file ./deployments/docker-compose.yml run mock_gen \
		--all \
		--dir=./pkg/service/ \
		--outpkg=mockservice \
		--output=./mocks/service \
		--keeptree;
list-service:
	docker-compose --file ./deployments/docker-compose.yml ps
exec-service:
	docker-compose --file ./deployments/docker-compose.yml exec $(service) bash
logs:
	docker-compose --file ./deployments/docker-compose.yml logs -f $(service)
destroy:
	docker-compose --file ./deployments/docker-compose.yml down --remove-orphans
