APP_NAME = my-app

.PHONY: build test deploy clean

build:
	docker build -t $(APP_NAME) .

clean:
	docker-compose down
	docker system prune -f

compose-up:
	docker-compose up -d

compose-down:
	docker-compose down

test: build
	docker run $(APP_NAME) pytest

test:
	docker-compose run web pytest

deploy: test
	docker push $(APP_NAME):latest

deploy:
	docker push my-registry/my-app:latest

all: build compose-up
