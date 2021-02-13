COMMIT    := $(shell git describe --no-match --always --dirty)
BRANCH    := $(shell git rev-parse --abbrev-ref HEAD)

REPO := github.com/lobshunter86/goproxy

GOENV   := GO111MODULE=on CGO_ENABLED=0
GO 		:= $(GOENV) go

LDFLAGS := -w -s
LDFLAGS += -X "$(REPO)/pkg/version.GitHash=$(COMMIT)"
LDFLAGS += -X "$(REPO)/pkg/version.GitBranch=$(BRANCH)"

default: lint server-side client-side

server-side: 
	$(GO) build -ldflags '$(LDFLAGS)' -o bin/server ./cmd/server/

client-side: 
	$(GO) build -ldflags '$(LDFLAGS)' -o bin/client ./cmd/client/

lint:
	golangci-lint run

test:
	$(GO) test ./... -coverpkg ./pkg/proxy  -covermode atomic -coverprofile bin/cover.out

clean:
	rm bin/*