RSCLI_VERSION=v0.0.26-alpha
rscli-version:
	@echo $(RSCLI_VERSION)

buildc:
	docker build -t matreshka-be:local .