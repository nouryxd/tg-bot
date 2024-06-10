dev:
	go build -o "tgbot.out" ./cmd/bot
	./tgbot.out
	
docker:
	docker compose build 
	docker compose up

dockerprod:
	docker compose build 
	docker compose up -d

