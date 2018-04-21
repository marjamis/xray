FROM golang:1.10.0 as building

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh && mkdir cfn-signaler

WORKDIR /go/src/application/

COPY Gopkg.lock .
COPY Gopkg.toml .
RUN dep ensure -vendor-only

COPY main.go /go/src/application/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:3.7
LABEL maintainer=marjamis@amazon.com \
      version=testing

RUN apk add --no-cache tini
ENTRYPOINT ["/sbin/tini", "--"]

USER nobody
COPY --from=building /go/src/application/app /
# Run your program under Tini
CMD ["/app"]
