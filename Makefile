deps:
	-rm -rf vendor
	-rm -f go.mod
	-rm -f go.sum
	GO111MODULE=on go mod init
	GO111MODULE=on go mod vendor
clean:
	go clean
	# docker rm $(docker ps -a -q)
	# docker rmi $(docker images -q)
test:
	go test ./...
install:
	make clean
	make deps
	go install
deploy:
	make clean
	make deps
	heroku container:login
	heroku container:push web -a tfclass
	heroku container:release web -a tfclass