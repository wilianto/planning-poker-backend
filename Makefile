generate:
	rm -rf ./model/schema/ent
	go generate ./model

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

generate-doc:
	swag init -q false