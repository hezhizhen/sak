PROJECT_PKG = github.com/hezhizhen/sak
CLI_EXE     = sak
CLI_PKG     = $(PROJECT_PKG)/cmd/sak
GIT_COMMIT  = $(shell git rev-parse HEAD)
GIT_TAG     = $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
GIT_DIRTY   = $(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")
LDFLAGS = -w
ifdef VERSION
	LDFLAGS += -X $(PROJECT_PKG)/pkg/version.BuildMetadata=$(VERSION)
endif
ifneq ($(GIT_TAG),)
	LDFLAGS += -X $(PROJECT_PKG)/pkg/version.BuildMetadata=
endif
LDFLAGS += -X $(PROJECT_PKG)/pkg/version.GitCommit=${GIT_COMMIT}
LDFLAGS += -X $(PROJECT_PKG)/pkg/version.GitTreeState=${GIT_DIRTY}

.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	git status --ignored --short | grep '^!! ' | sed 's/!! //' | xargs rm -rf

.PHONY: build-mac
build-mac:
	GOARCH=amd64 CGO_ENABLED=0 GOOS=darwin go build -v --ldflags="$(LDFLAGS)" \
		-o bin/$(CLI_EXE) $(CLI_PKG)

.PHONY: build-linux
build-linux:
	GOARCH=amd64 CGO_ENABLED=0 GOOS=linux go build -v --ldflags="$(LDFLAGS)" \
		-o bin/$(CLI_EXE) $(CLI_PKG)

update-local:
	go install ./...
	sak completion fish > ~/.config/fish/completions/sak.fish