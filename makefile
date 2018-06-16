unit:
	go test ./...
build:
	go build -o cli-mgr main.go
	chmod +x cli-mgr