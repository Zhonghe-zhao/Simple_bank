DB_URL=postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable

network:
	docker network create bank_network

postgres:
	docker run --name postgres12 --network bank_network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	 docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	 docker exec -it postgres12 dropdb  simple_bank

migrateup:
	migrate -path ./db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path ./db/migration -database "$(DB_URL)" -verbose down

migrateup1:
	migrate -path ./db/migration -database "$(DB_URL)" -verbose up 1

migratedown1:
	migrate -path ./db/migration -database "$(DB_URL)" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover 
	
server:
	go run main.go

mock:
	mockgen -destination db/mock/store.go project/simplebank/db/sqlc Store 

db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

proto:
	del /Q pb\*.go\
	del /Q doc\swagger\*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
    proto/*.proto
	statik -src=./doc/swagger -dest=./doc -f

evans:
	evans --host localhost --port 9090 -r repl

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server mock migrateup1 migratedown1 db_docs db_schema proto