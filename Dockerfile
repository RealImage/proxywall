FROM golang:1.9
RUN go get github.com/RealImage/proxywall
RUN go install github.com/RealImage/proxywall
ENTRYPOINT /go/bin/proxywall