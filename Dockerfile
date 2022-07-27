FROM golang:1.18-alpine

# copy source code into the $GOPATH and switch to that directory
COPY  . ${GOPATH}/src/github.com/bigbes/go-echo
WORKDIR ${GOPATH}/src/github.com/bigbes/go-echo

# compile source code and copy into $PATH
RUN go install -mod=vendor ./cmd/go-echo/main.go && cp ${GOPATH}/bin/main /bin/go-echo

# the default command runs the service in the foreground
CMD ["/bin/go-echo"]
EXPOSE 9000
EXPOSE 9090