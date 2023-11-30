build:
	go build -o bin/stocks-data-aggregator

run: build
	./bin/stocks-data-aggregator

test:
	go test -v ./...

deploy: 
	docker build . -t "stocks-data-aggregator" && docker-compose --env-file .env up -d 

stop:
	docker-compose down

clean: stop
	docker system prune -f && docker volume prune -af && powershell rm -r -force ./volumes/