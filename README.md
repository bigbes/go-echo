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
