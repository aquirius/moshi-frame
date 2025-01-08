clean:
	docker-compose -p sprout-backend down && docker-compose -p sprout-backend_redisfssprout_1 down
	docker-compose -p sprout-backend up -d && docker-compose -p  sprout-backend_redisfssprout_1 up -d

run:
	SPROUT_MYSQL="root:moshi@tcp(127.0.0.1:3311)/sprout?collation=utf8_general_ci" \
	go run ./main.go

test-user:
	go test ./internal/systems/user


test-all:
	go test -v ./...

load-schemas:
	go run ./schemas/load.go