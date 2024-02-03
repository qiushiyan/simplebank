.PHONY: createdb dropdb

KIND_CLUSTER    := qs-starter-cluster
KIND            := kindest/node:v1.29.0@sha256:eaa1450915475849a73a9227b8f201df25e55e268e5d619312131292e324d570
SERVICE_NAME    := bank-api
VERSION         := 0.0.1
SERVICE_IMAGE   := $(SERVICE_NAME):$(VERSION)

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


dev-up:
	kind create cluster \
		--image $(KIND) \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/dev/kind-config.yaml

	kubectl wait --timeout=120s --namespace=local-path-storage --for=condition=Available deployment/local-path-provisioner


dev-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

dev-down:
	kind delete cluster --name $(KIND_CLUSTER)

bank-api:
	docker build \
		-f zarf/docker/dockerfile.bank-api \
		-t $(SERVICE_IMAGE) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

generate:
	sqlc generate

test:
	go test -v -cover ./...

check:
	nilaway app/*