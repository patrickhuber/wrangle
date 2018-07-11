unit:
	go test ./...

build:
	go build -o bin/wrangle main.go

restore:
	dep ensure
	
integration:
	go test ./... -tags=integration