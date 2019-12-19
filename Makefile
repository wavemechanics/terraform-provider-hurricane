.PHONY: help

# self-documenting Makefile thanks to
# https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

PLUGIN=terraform-provider-hurricane
PLUGIN_DIR=$(HOME)/.terraform.d/plugins

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

install: build  ## install versioned binary
	mkdir -p $(PLUGIN_DIR)
	export SEMVER=`tools/get-binary-semver.sh $(PLUGIN)` && \
	export DEST="$(PLUGIN_DIR)/$(PLUGIN)_$$SEMVER" && \
	cp -a $(PLUGIN) $$DEST && \
	echo "To uninstall, remove $$DEST"

build:	## Compile the provider
	tools/brand.sh config/version.go `tools/get-git-version.sh`
	go build ./cmd/...

clean::	## remove compiled binary
	rm -f $(PLUGIN)

release: validate	## Tag this version as a release
	@test -n "$(SEMVER)" || ( echo "must set SEMVER for release"; false)
	tools/brand.sh config/version.go "$(SEMVER)"
	go build ./cmd/...
	test -d .git || (echo "must be in git working directory"; false)
	git add config/version.go
	git commit -m "branded $(SEMVER)"
	git push
	git tag "$(SEMVER)"
	git push --tag

validate:	## validate the CircleCI configuration
	circleci config validate
