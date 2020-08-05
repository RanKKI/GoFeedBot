FROM golang:1.14

WORKDIR /go/src/app
COPY . .

ENV token EMPTY
ENV debug 0
ENV proxy EMPTY
ENV interval 30m

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["app"]