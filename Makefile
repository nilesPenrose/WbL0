postgres-build:
	docker-compose up -d --build postgres

nats-build:
	docker-compose up -d --build nats

postgres-up:
	docker-compose up -d postgres

nats-up:
	docker-compose up -d nats

service:
	go run cmd/main.go

front:
	simple_web_interface/venv/bin/python3 simple_web_interface/main.py

create-venv:
	python3 -m venv simple_web_interface/venv && source simple_web_interface/venv/bin/activate && pip install -r requirements.txt

.PHONY:all
all: postgres-build nats-build service front
