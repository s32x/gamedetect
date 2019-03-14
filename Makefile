deps:
	-rm -rf vendor
	-rm -f go.mod
	-rm -f go.sum
	go clean
	GO111MODULE=on go mod init
	GO111MODULE=on go mod vendor
test:
	go test ./...
install:
	make deps
	go install
deploy:
	make deps
	heroku container:login
	heroku container:push web -a gamedetect
	heroku container:release web -a gamedetect