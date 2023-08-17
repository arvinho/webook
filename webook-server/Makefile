.PHONY: docker
docker:
	@rm webook || true
	@GOOS=linux GOARCH=amd64 go build -o webook .
	@docker rmi -f arvin/webook:v0.0.1
	@docker build -t arvin/webook:v0.0.1 .