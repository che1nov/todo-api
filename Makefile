BINARY=todo-api
DOCKER_IMAGE=todo-api
DB_FILE=./todo.db

build:
	go build -o $(BINARY) ./cmd/main.go

run: build
	./$(BINARY)

clean:
	rm -f $(BINARY)
	rm -f $(DB_FILE)

docker-build:
	docker build -t $(DOCKER_IMAGE) .

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down

docker-clean:
	docker rmi $(DOCKER_IMAGE) || true