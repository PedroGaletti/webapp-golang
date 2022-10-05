install:
	@go get
	@go mod tidy
	@go mod vendor

build:
	@go build

run:
	@go run main.go

clean:
	@rm -rf ./vendor
	@rm -rf ./go.sum
	@rm -rf ./webapp