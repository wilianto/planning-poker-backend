run:
	cp .env.example .env
	docker compose --env-file .env up -d
	rm -rf ./model/schema/ent
	go generate ./model
	go run main.go

clean:
	docker compose --env-file .env down -v
	rm -rf ./model/schema/ent
	go clean

run-test:
	cp .env.example .env
	docker compose --env-file .env up -d
	rm -rf ./model/schema/ent
	go generate ./model
	go test -count=1 -v ./...

generate-model:
	rm -rf ./model/schema/ent
	go generate ./model

generate-doc:
	swag init -q false