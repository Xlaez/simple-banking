postgres: 
	sudo docker run --name postgres13-old -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres\:13-alpine
createdb:
	sudo docker exec -it postgres13-old createdb --username=root --owner=root simple_bank
dropdb:
	sudo docker exec -it postgres13-old dropdb simple_bank
migrateup: 
	migrate -path ./db/migration/ -database "postgresql://root:password@localhost:5432/simple_bank?sslmode=disable" -verbose up
migratedown:
	migrate -path ./db/migration/ -database "postgresql://root:password@localhost:5432/simple_bank?sslmode=disable" -verbose down
test:
	go test -v -cover ./...
sqlc:
	sqlc generate
.PHONY: postgres createdb dropdb migrateup migratedown sqlc test