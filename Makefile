setup:
	pip install pre-commit
	pre-commit install

reinstall-hooks:
	pre-commit clean
	pre-commit install

check:
	pre-commit run --all-files

build:
	go build ./services/streaming-service/...
	go build ./services/auth-service/...
	go build ./services/chat-service/...

docker-up:
	docker-compose up --build -d

docker-down:
	docker-compose down

docker-clean:
	docker-compose down --volumes --rmi all
