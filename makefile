#Down container
down:
	docker-compose down
#Up containner
up: down
	docker-compose up -d
#Formater
fmt:
	go fmt ./...
#Linter
lint:
	docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v1.50.1 golangci-lint run
#Test
test:
	go test -count=1 -p=1 ./... -v | grep -v "no test"