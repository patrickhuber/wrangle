unit:
	ginkgo -p -r -race -randomizeAllSpecs -randomizeSuites -skipPackage vendor .

build:
	go build -o bin/wrangle main.go

restore:
	dep ensure
	
integration:
	go test ./... -tags=integration