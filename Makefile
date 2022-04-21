#Variables
db_name := simple_bank
db_username := root
db_password := secret123
container_name := sbpostgres
command := version

#Targets
.PHONY: setup
setup: start-postgres wait create-db migrate-up

.PHONY: shutdown
shutdown: remove-postgres

.PHONY: wait
wait:
	sleep 3

.PHONY: remove-postgres
remove-postgres:
	docker stop $(container_name) && docker rm $(container_name)

.PHONY: start-postgres
start-postgres:
	docker run --name $(container_name) -p 5432:5432 -e POSTGRES_USER=$(db_username) -e POSTGRES_PASSWORD=$(db_password) -d postgres

.PHONY: create-db
create-db:
	docker exec -it $(container_name) createdb --username=$(db_username) --owner=$(db_username) $(db_name)

.PHONY: drop-db
drop-db:
	docker exec -it $(container_name) dropdb $(db_name)

.PHONY: migrate-up
migrate-up:
	migrate -path db/sqlc/migrations -database "postgresql://root:secret123@localhost:5432/simple_bank?sslmode=disable" -verbose up

.PHONY: migrate-down
migrate-down:
	migrate -path db/sqlc/migrations -database "postgresql://$(db_username):$(db_password)@localhost:5432/$(db_name)?sslmode=disable" -verbose down

.PHONY: sqlc
sqlc:
	docker run --rm -v $(shell pwd):/src -w /src kjconroy/sqlc generate

.PHONY: test
test:
	go test -cover ./...


