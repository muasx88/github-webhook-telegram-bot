BINARY=github-webhook-telegram-bot

docker-up:
	@ docker-compose up --build

dev:
	go run ./app/main.go

build:
	@ printf "Build... "

	@ go build \
 		-trimpath \
 		-o target/${BINARY} \
 		./app/

	@ echo "done"

.PHONY: docker-up dev build