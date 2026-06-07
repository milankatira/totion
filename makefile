.PHONY: dev build run clean

dev:
	go run .

build:
	@go build -o totion .

run: build
	@./totion

clean:
	rm -f totion
