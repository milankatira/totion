.PHONY: dev build run clean

dev:
	go run main.go

build:
	@go build -o totion .

run: build
	@./totion

clean:
	rm -f totion
