unit:
	ginkgo -p -r -race -randomizeAllSpecs -randomizeSuites .

build:
	go build -o bin/wrangle

restore:
	go mod tidy
	
integration:
	go test ./... -tags=integration