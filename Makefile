postgres: 
	docker run --name postgres13-old -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres\:13-alpine
createdb:
	docker exec -it postgres13-old createdb --username=root --owner=root bank
dropdb:
	docker exec -it postgres13-old dropdb bank
migrateup: 
	migrate -path ./db/migration/ -database "postgresql://root:password@localhost:5432/bank?sslmode=disable" -verbose up
migrateup1: 
	migrate -path db/migration/ -database "postgresql://root:password@localhost:5432/bank?sslmode=disable" -verbose up 1
migratedown:
	migrate -path ./db/migration/ -database "postgresql://root:password@localhost:5432/bank?sslmode=disable" -verbose down
migratedown1:
	migrate -path db/migration/ -database "postgresql://root:password@localhost:5432/bank?sslmode=disable" -verbose down 1	
test:
	go test -v -cover ./...
sqlc:
	sqlc generate
serve:
	go run main.go
migrateupdate:
	migrate create -ext sql -dir db/migration -seq add_users
opendb:
	docker exec -it postgres13-old psql -U root -d bank
.PHONY: 
	postgres createdb dropdb migrateup migratedown sqlc test serve migrateupdate cleandb migratedown1 migrateup1