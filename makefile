clean:
	@rm -rf bin


build:
	clear
	@go build -o bin/timerino


run: build
	@clear
	@go build -o bin/timerino
	@./bin/timerino 0 0 5 "This is the countdown message!"


test:
	@clear
	@go test
