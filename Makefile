GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAGS=-s -w

cli:
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/log-groups cmd/log-groups/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/log-group-streams cmd/log-group-streams/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/log-stream-events cmd/log-stream-events/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/empty-streams cmd/empty-streams/main.go
