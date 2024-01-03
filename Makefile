postgres:
	docker run --name flex_admin -p 5105:5432 -e POSTGRES_PASSWORD=Si73gangan -d postgres:12-alpine
createdb:
	docker exec -it flex_admin createdb --username=postgres --owner=postgres flex_server

dropdb:
	docker exec -it flex_admin dropdb flex_server
	
createmigrate:
	migrate create -ext sql -dir db/migration -seq init_scheme
migrateup:
	docker run -v /Users/uwa/Documents/flex_server/db/migration:/migrations --network host migrate/migrate -path=/migrations/ -database "postgresql://postgres:Si73gangan@0.0.0.0:5105/flex_server?sslmode=disable" up

migratedown:
	docker run -v /Users/uwa/Documents/flex_server/db/migration:/migrations --network host migrate/migrate -path=/migrations/ -database "postgresql://postgres:Si73gangan@0.0.0.0:5105/flex_server?sslmode=disable" down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...	

proto:
	rm -f pb/*.go
	protoc --plugin=protoc-gen-go=/Users/uwa/go/bin/protoc-gen-go --plugin=protoc-gen-go-grpc=/Users/uwa/go/bin/protoc-gen-go-grpc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    proto/*.proto

evans:
	evans --host localhost --port 9090  -r repl

server:
	go run main.go	
.PHONY: postgres createdb dropdb migrateup migratedown sqlc server proto