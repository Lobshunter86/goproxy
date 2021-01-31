COMMIT    := $(shell git describe --no-match --always --dirty)
BRANCH    := $(shell git rev-parse --abbrev-ref HEAD)

REPO := github.com/lobshunter86/goproxy

LDFLAGS := -w -s
LDFLAGS += -X "$(REPO)/pkg/version.GitHash=$(COMMIT)"
LDFLAGS += -X "$(REPO)/pkg/version.GitBranch=$(BRANCH)"

default: vet server-side client-side

server-side: 
	go build -ldflags '$(LDFLAGS)' -o bin/server ./cmd/server/

client-side: 
	go build -ldflags '$(LDFLAGS)' -o bin/client ./cmd/client/

vet:
	go vet ./...