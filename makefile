unit:
	go test ./...

build:
	go build -o cli-mgr main.go

restore:
	dep ensure
	
integration:
	go test ./... -tags=integration