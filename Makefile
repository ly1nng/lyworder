BINARY_NAME=lyworder

docker-run:
	docker compose up -d

docker-image:
	docker build -t ${BINARY_NAME}  . 

build:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o bin/${BINARY_NAME} main.go
run:
	./${BINARY_NAME}

build_and_run:build run
clean:
	go clean
	rm bin/${BINARY_NAME}
