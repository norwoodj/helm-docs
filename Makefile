build:
	go build -ldflags '-extldflags "-static"'

clean:
	go clean
