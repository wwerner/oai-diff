
FROM golang:1.11.2 as build-image
RUN mkdir -p $GOPATH/src/github.com/wwerner/oaidiff
ADD . $GOPATH/src/github.com/wwerner/oaidiff
WORKDIR $GOPATH/src/github.com/wwerner/oaidiff

RUN go get -u -v github.com/golang/dep/cmd/dep
RUN dep ensure -vendor-only -v

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o oaidiff main.go
RUN cp oaidiff /


FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=build-image /oaidiff /home/

WORKDIR /home

ENTRYPOINT ["./oaidiff"]
CMD []
