cli:
	go build -mod vendor -ldflags="-s -w" -o bin/log-groups cmd/log-groups/main.go
	go build -mod vendor -ldflags="-s -w" -o bin/log-group-streams cmd/log-group-streams/main.go
	go build -mod vendor -ldflags="-s -w" -o bin/log-stream-events cmd/log-stream-events/main.go
	go build -mod vendor -ldflags="-s -w" -o bin/empty-streams cmd/empty-streams/main.go
