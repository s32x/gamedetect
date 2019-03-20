deps:
	-rm -rf vendor
	-rm -f go.mod
	-rm -f go.sum
	-rm -rf service/static/test/.DS_Store
	go clean
	GO111MODULE=on go mod init
	GO111MODULE=on go mod vendor
test:
	go test ./...
install:
	make deps
	go install
run:
	DEMO=true go run main.go