# ==============================================================================
# Dependencies
KIND_CLUSTER    := bank-system-cluster
NAMESPACE       := simplebank
KIND            := kindest/node:v1.29.0@sha256:eaa1450915475849a73a9227b8f201df25e55e268e5d619312131292e324d570
BASE_IMAGE_NAME := qiushiyan/simplebank
SERVICE_NAME    := bank-api
APP             := bank-api
VERSION         := 0.0.1
SERVICE_IMAGE   := $(BASE_IMAGE_NAME)/$(SERVICE_NAME):$(VERSION)

# ==============================================================================
# Database
pg:
	docker run --name postgres -e POSTGRES_PASSWORD=postgres -p 5433:5432 -d bank-api-postgres

pgcli:
	docker exec -it bank-api-postgres psql -U postgres

createdb:
	docker exec -it bank-api-postgres createdb --username=postgres --owner=postgres bank

dropdb:
	docker exec -it bank-api-postgres dropdb bank --username=postgres

migrate-create:
	migrate create -ext sql -dir business/db/migrations -seq init_schema

migrate-up:
	migrate -path business/db/migrations -database "postgresql://postgres:postgres@localhost:5433/bank?sslmode=disable" --verbose up

migrate-down:
	migrate -path business/db/migrations -database "postgresql://postgres:postgres@localhost:5433/bank?sslmode=disable" --verbose down

generate:
	sqlc generate

# ==============================================================================
# Running locally
run-local:
	go run app/services/bank-api/main.go | go run app/tooling/logfmt/main.go --service=$(SERVICE_NAME)

run-local-help:
	go run app/services/bank-api/main.go -h

# ==============================================================================
# Running from within k8s/kind
dev-up:
	kind create cluster \
		--image $(KIND) \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/dev/kind-config.yaml

	kubectl wait --timeout=120s --namespace=local-path-storage --for=condition=Available deployment/local-path-provisioner

dev-load:
	kind load docker-image $(SERVICE_IMAGE) --name $(KIND_CLUSTER)

dev-apply:
	kustomize build zarf/k8s/dev/bank-api | kubectl apply -f -
	kubectl wait pods --namespace=$(NAMESPACE) --selector app=$(APP) --timeout=120s --for=condition=Ready

dev-down:
	kind delete cluster --name $(KIND_CLUSTER)

dev-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

dev-describe-deployment:
	kubectl describe deployment $(APP) --namespace=$(NAMESPACE)

dev-describe-bank-api:
	kubectl describe pod $(APP) --namespace=$(NAMESPACE)

dev-logs:
	kubectl logs --namespace=$(NAMESPACE) -l app=$(APP) --all-containers=true -f --tail=100 | go run app/tooling/logfmt/main.go --service=$(SERVICE_NAME)

dev-start: dev-up dev-load dev-apply

dev-restart:
	kubectl rollout restart deployment $(APP) --namespace=$(NAMESPACE)

dev-update: all dev-load dev-restart

dev-update-apply: all dev-load dev-apply


# ==============================================================================
# Building containers
all: bank-api

bank-api:
	docker build \
		-f zarf/docker/dockerfile.bank-api \
		-t $(SERVICE_IMAGE) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

# ==============================================================================
# Chores
tidy:
	go mod tidy

test:
	go test -v -cover ./...

check:
	nilaway ./app/*/**
