.PHONY: run dev build

run:
	go run main.go

dev:
	air

build:
	go build -o go-todo-sample .
