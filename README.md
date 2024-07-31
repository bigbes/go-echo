go-echo
=======

[![Docker Pulls](https://img.shields.io/docker/pulls/paddycarey/go-echo.svg)](https://hub.docker.com/r/paddycarey/go-echo/)

go-echo is a small HTTP server that echos the request headers and body back at a client as JSON in the response body.

http/2
------

go-echo supports HTTP/2 over cleartext TCP (h2c).
           
### Examples

Make HTTP/1.1 GET and/or POST request with upgrade to h2c. This is a common way for case when the server supports h2c but the client does not know that.    

```bash
curl --http2 -v 127.0.0.1:9001/get
curl --http2 -v 127.0.0.1:9001/post --data '{"username":"xyz","password":"xyz"}'
```

Make HTTP/2 GET and/or POST request with prior knowledge of h2c. This is a common way for case when the client knows that the server supports h2c.
            
```bash
curl --http2 --http2-prior-knowledge -v 127.0.0.1:9001/get
curl --http2 --http2-prior-knowledge -v 127.0.0.1:9001/post --header "Content-Type: application/json" --data '{"username":"xyz","password":"xyz"}'
```

gRPC
------

This is a simple gRPC server with reflection support.

The request message containing the user's name.
Receives a request:
```
message HelloRequest {
  string name = 1;
}
```

The response message containing the greetings.
Replies with a response:
```
message HelloReply {
  string message = 1;
}
```

Download protobuf spec from:
https://github.com/grpc/grpc/blob/master/examples/protos/helloworld.proto


### Examples

For create simple request you can use:
```
grpcurl --plaintext --import-path . --proto ./helloworld.proto -d '{"name": "Test"}' ${GRPCSERVER_IP}:10050 helloworld.Greeter/SayHello
```

You can also build this server, run it with Kubernetes, run an Ubuntu container, and submit the previous request for testing.

```
docker build . -f ./Dockerfile --build-arg build_version="test" -t grpcserver:local

kubectl -n default run grpcserver --restart=Never --image=grpcserver:local --port=10050
```

You can use script for run:
```
#!/bin/sh

repeat_count=1200
i=1

while [ $i -le $repeat_count ]
do
    echo "Start $i"
    grpcurl --plaintext --import-path . --proto ./helloworld.proto -d '{"name": "Test"}' ${GRPCSERVER_IP}:10050 helloworld.Greeter/SayHello
    sleep 3
    i=$((i+1))
done

echo "All $repeat_count finished"
```


