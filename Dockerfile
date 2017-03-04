# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/wobeproject

# RUN export config_path="/go/src/github.com/wobeproject/config/config.toml"

RUN cd src/github.com/wobeproject && ./testscript.sh

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
# RUN go build main.go -env prod && go install main

RUN go install ./...
# RUN find / -name "config.toml"
# Run the outyet command by default when the container starts.
ENTRYPOINT ["/go/bin/wobeproject","--config","/go/src/github.com/wobeproject/config/config.toml"]

# Document that the service listens on port 8080.
EXPOSE 8081
