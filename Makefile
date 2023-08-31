.PHONY: build run shutdown

build:
	docker build --tag dnevsky/avito-trainee .

run:
	docker-compose up -d avito-trainee

shutdown:
	docker-compose down