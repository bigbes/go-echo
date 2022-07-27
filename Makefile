build:
	go build -o .artifacts/echo $(realpath ./cmd/echo)/*.go

image:
	docker build -t bigbes/go-echo .

image-push:
	docker image push bigbes/go-echo
