ALPINE          := alpine:3.18
KIND_CLUSTER    := simplebank-cluster
NAMESPACE       := simplebank
KIND            := kindest/node:v1.29.0@sha256:eaa1450915475849a73a9227b8f201df25e55e268e5d619312131292e324d570
TELEPRESENCE    := datawire/tel2:2.13.1
BASE_IMAGE_NAME := qiushiyan/simplebank
SERVICE_NAME    := bank-api
APP             := bank-api
VERSION         := 0.0.1
SERVICE_IMAGE   := $(BASE_IMAGE_NAME)/$(SERVICE_NAME):$(VERSION)
CONTAINER_NAME  := bank-api-postgres
POSTGRES        := postgres:latest
DATABASE_URL 	:= postgres://postgres:postgres@localhost:5432/bank?sslmode=disable

compose:
	docker-compose up -d

# ==============================================================================
# Install dependencies
dev-gotooling:
	go install github.com/divan/expvarmon@latest
	go install github.com/rakyll/hey@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install go.uber.org/mock/mockgen@latest
	go install github.com/swaggo/swag/cmd/swag@latest


dev-brew:
	brew update
	brew list kind || brew install kind
	brew list kubectl || brew install kubectl
	brew list kustomize || brew install kustomize
	brew list pgcli || brew install pgcli


dev-docker:
	docker pull $(POSTGRES)

# ==============================================================================
# Database
pg:
	docker run --name $(CONTAINER_NAME) -e POSTGRES_PASSWORD=postgres -p 5433:5432 -d postgres

pgcli:
	docker exec -it $(CONTAINER_NAME) psql -U postgres

createdb:
	docker exec -it $(CONTAINER_NAME) createdb --username=postgres --owner=postgres bank

dropdb:
	docker exec -it $(CONTAINER_NAME) dropdb bank --username=postgres

migrate-create:
	migrate create -ext sql -dir business/db/migration -seq init_schema

migrate-up:
	migrate -path business/db/migration -database $(DATABASE_URL) --verbose up

migrate-down:
	migrate -path business/db/migration -database $(DATABASE_URL) --verbose down

migrate-up-test:
	migrate -path business/db/migration -database "postgresql://postgres:postgres@localhost:5432/bank_test?sslmode=disable" --verbose up

generate:
	sqlc generate
	mockgen -destination=./business/db/mock/mockstore.go -package=mockdb github.com/qiushiyan/simplebank/business/db/core Store

# ==============================================================================
# Running locally
run-local:
	go run app/services/bank-api/main.go | go run app/tooling/logfmt/main.go --service=$(SERVICE_NAME)

run-local-help:
	go run app/services/bank-api/main.go -h

# ==============================================================================
# Running from within k8s/kind
dev-start: dev-up dev-load dev-apply

dev-up:
	kind create cluster \
		--image $(KIND) \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/dev/kind-config.yaml

	kubectl wait --timeout=120s --namespace=local-path-storage --for=condition=Available deployment/local-path-provisioner
	kind load docker-image $(POSTGRES) --name $(KIND_CLUSTER)

dev-load:
	kind load docker-image $(SERVICE_IMAGE) --name $(KIND_CLUSTER)

dev-apply:
	# database
	kustomize build zarf/k8s/dev/database | kubectl apply -f -
	kubectl rollout status --namespace=$(NAMESPACE) --watch --timeout=120s sts/database

	# bank-api
	kustomize build zarf/k8s/dev/bank-api | kubectl apply -f -
	kubectl wait pods --namespace=$(NAMESPACE) --selector app=$(APP) --timeout=120s --for=condition=Ready

dev-forward:
	kubectl port-forward svc/bank-api-svc 3000:3000 4000:4000 -n $(NAMESPACE)

dev-down:
	kind delete cluster --name $(KIND_CLUSTER)

dev-status:
	kubectl get nodes -o wide -n $(NAMESPACE)
	kubectl get svc -o wide -n $(NAMESPACE)
	kubectl get pods -o wide -n $(NAMESPACE)

# make dev-status | go run app/tooling/k8sfmt/main.go

dev-describe-deployment:
	kubectl describe deployment $(APP) --namespace=$(NAMESPACE)

dev-describe-pod:
	kubectl describe pod $(APP) --namespace=$(NAMESPACE)

dev-describe-svc:
	kubectl describe svc $(APP) --namespace=$(NAMESPACE)

dev-describe-gateway:
	kubectl describe gateway $(APP) --namespace=$(NAMESPACE)

dev-logs:
	kubectl logs --namespace=$(NAMESPACE) -l app=$(APP) --all-containers=true -f --tail=100 | go run app/tooling/logfmt/main.go --service=$(SERVICE_NAME)

dev-logs-db:
	kubectl logs --namespace=$(NAMESPACE) -l app=database --all-containers=true -f --tail=100

dev-logs-init:
	kubectl logs --namespace=$(NAMESPACE) -l app=$(APP) -f --tail=100 -c init-migrate


dev-restart:
	kubectl rollout restart deployment $(APP) --namespace=$(NAMESPACE)

dev-update: all dev-load dev-restart

dev-update-apply: all dev-load dev-apply


# ==============================================================================
# Building containers
all: bank-api

frontend:
	docker build -f zarf/docker/dockerfile.frontend -t simplebank-frontend .

bank-api:
	docker build \
		-f zarf/docker/dockerfile.bank-api \
		-t $(SERVICE_IMAGE) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

# ==============================================================================
# Chores
conf-help:
	go run app/services/bank-api/main.go -h

docs:
	swag init --dir app/services/bank-api -o app/services/bank-api/docs --parseDependency --parseInternal --parseDepth 1
	swag fmt

tidy:
	go mod tidy

test:
	go test -v -cover ./...

check:
	nilaway ./app/*/**

metrics-view:
	expvarmon -ports="$(SERVICE_NAME).$(NAMESPACE).svc.cluster.local:4000" -vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"

metrics-view-local:
	expvarmon -ports="localhost:4000" -vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"


