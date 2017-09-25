VERSION=$(shell cat VERSION)
IMAGE_NAME=ibigbug/ss-account

all:build push

build:
	docker build -t $(IMAGE_NAME):$(VERSION) .

push:
	docker push $(IMAGE_NAME):$(VERSION)

