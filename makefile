unit:
	go test ./...

build:
	go build -o cli-mgr main.go
	
integration:
	go test ./... -tags=integration