

### General Rscli info
RSCLI_VERSION=v0.0.31
rscli-version:
	@echo $(RSCLI_VERSION)

build-local-container:
	docker buildx build \
			--load \
			--platform linux/arm64 \
			-t red_cart:local .