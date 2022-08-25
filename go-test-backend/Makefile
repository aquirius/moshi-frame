clean:
	docker-compose -p go-test-backend down && docker-compose -p go-test-backend_redisfsmilon_1 down
	docker-compose -p go-test-backend up -d && docker-compose -p  go-test-backend_redisfsmilon_1 up -d

run:
	MILON_MYSQL="milon:milon@tcp(127.0.0.1:3311)/milon?collation=utf8_general_ci" \
	go run ./main.go

test-user:
	go test ./internal/systems/user


test-all:
	go test -v ./...

load-schemas:
	go run ./schemas/load.go