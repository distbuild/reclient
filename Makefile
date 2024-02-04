# Build

.PHONY: FORCE

build: rc-build
.PHONY: build

install: rc-install
.PHONY: install

# Non-PHONY targets (real files)

rc-build: FORCE
	./scripts/build.sh

rc-install: FORCE
	./scripts/install.sh
