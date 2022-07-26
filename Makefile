REGISTRY_NAME ?= stardotstar
IMAGE_VERSION ?= 0.0.3
IMAGE_NAME ?= aurora-tracker
IMAGE_TAG ?= ${REGISTRY_NAME}/${IMAGE_NAME}:${IMAGE_VERSION}
OUTPUT_TYPE ?= registry
BUILDX_BUILDER_NAME ?= img-builder

ARCH ?= amd64
ALL_ARCH.linux = amd64 arm64
ALL_OS_ARCH.linux = $(foreach arch, ${ALL_ARCH.linux}, linux-$(arch))

.PHONY: build
build:
	CGO_ENABLED=0 GOARCH=${ARCH} GOOS=linux go build -a -o _output/${ARCH}/aurora-tracker .

.PHONY: docker-buildx-builder
docker-buildx-builder:
	@if ! docker buildx ls | grep $(BUILDX_BUILDER_NAME); then \
		docker buildx create --name $(BUILDX_BUILDER_NAME) --use; \
		docker buildx inspect $(BUILDX_BUILDER_NAME) --bootstrap; \
	fi

.PHONY: container-all
container-all:
	for arch in $(ALL_ARCH.linux); do \
		ARCH=$${arch} $(MAKE) build; \
		ARCH=$${arch} $(MAKE) container-linux; \
	done

.PHONY: container-linux
container-linux: docker-buildx-builder
	docker buildx build \
		--no-cache \
		--output=type=$(OUTPUT_TYPE) \
		--platform="linux/$(ARCH)" \
		-t $(IMAGE_TAG)-linux-$(ARCH) -f Dockerfile .

.PHONY: push-manifest
push-manifest:
	docker manifest create --amend $(IMAGE_TAG) $(foreach osarch, $(ALL_OS_ARCH.linux), $(IMAGE_TAG)-${osarch})
	docker manifest push --purge $(IMAGE_TAG)
	docker manifest inspect $(IMAGE_TAG)