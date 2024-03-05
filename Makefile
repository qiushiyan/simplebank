ALPINE          := alpine:3.18
KIND_CLUSTER    := bank-system-cluster
NAMESPACE       := simplebank
KIND            := kindest/node:v1.29.0@sha256:eaa1450915475849a73a9227b8f201df25e55e268e5d619312131292e324d570
TELEPRESENCE    := datawire/tel2:2.13.1
BASE_IMAGE_NAME := qiushiyan/simplebank
SERVICE_NAME    := bank-api
APP             := bank-api
VERSION         := 0.0.1
SERVICE_IMAGE   := $(BASE_IMAGE_NAME)/$(SERVICE_NAME):$(VERSION)
CONTAINER_NAME  := bank-api-postgres

# ==============================================================================
# Install dependencies
dev-gotooling:
	go install github.com/divan/expvarmon@latest
	go install github.com/rakyll/hey@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install go.uber.org/mock/mockgen@latest

dev-brew:
	brew update
	brew list kind || brew install kind
	brew list kubectl || brew install kubectl
	brew list kustomize || brew install kustomize
	brew list pgcli || brew install pgcli

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
	migrate -path business/db/migration -database "postgresql://postgres:postgres@localhost:5433/bank?sslmode=disable" --verbose up

migrate-down:
	migrate -path business/db/migration -database "postgresql://postgres:postgres@localhost:5433/bank?sslmode=disable" --verbose down

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
dev-start: dev-up dev-load dev-apply

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
conf-help:
	go run app/services/bank-api/main.go -h

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

gen-mock:
	mockgen -destination=./business/db/mock/mockstore.go -package=mockdb github.com/qiushiyan/simplebank/business/db/core Store