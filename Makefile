PROJECT_NAME=gosocket
BUILD_VERSION=1.0.0

DOCKER_IMAGE=$(PROJECT_NAME):$(BUILD_VERSION)
GO_BUILD_ENV=CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on

compose_dev: docker
	docker tag gosocket:1.0.0 ngxvu/gosocket:1.0.0; \
	docker push ngxvu/gosocket:1.0.0;

#B1
build:
	$(GO_BUILD_ENV) go build -v -o $(PROJECT_NAME)-$(BUILD_VERSION).bin main.go

#B2
docker_prebuild: build
	mv $(PROJECT_NAME)-$(BUILD_VERSION).bin deploy/$(PROJECT_NAME).bin; \

#B4
docker_build:
	cd deploy; \
	docker build --rm -t $(DOCKER_IMAGE) .;

#B5
docker_postbuild:
	cd deploy; \
	rm -rf $(PROJECT_NAME).bin 2> /dev/null;\

#B0
docker: docker_prebuild docker_build docker_postbuild
