APP_NAME = go_api


build:
	@go build -o ./bin/$(APP_NAME)
run: build
	@./bin/$(APP_NAME)
test:
	@go test -v ./... 

# add a seed @go build -o ./bin/$(APP_NAME) --seed true script
seed: build
	@./bin/$(APP_NAME) --seed
