init:
	-rm -rf ./vendor go.mod go.sum
	GO111MODULE=on go mod init

deps:
	-rm -rf ./vendor go.sum
	GO111MODULE=on go mod vendor
	
test:
	go test ./...

run:
	DEMO=true go run main.go