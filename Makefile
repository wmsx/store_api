.PHONY: docker
docker:
	docker build . -t store-api:latest
