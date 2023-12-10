build:
	templ generate && go build -o bin/cfd

run: build
	./bin/cfd

test:
	go test -v ./...

deploy: 
	templ generate && docker build . -t "cfd" && docker-compose --env-file .env up -d 

stop:
	docker-compose down

clean: stop
	docker system prune -f && docker volume prune -af && powershell rm -r -force ./volumes/