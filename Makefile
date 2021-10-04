cli:
	go build -mod vendor -o bin/log-groups cmd/log-groups/main.go
	go build -mod vendor -o bin/log-group-streams cmd/log-group-streams/main.go
