FROM golang:1.9 as build
RUN go get -u github.com/Masterminds/glide
COPY glide.* ./
RUN glide install
COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .


FROM alpine:latest
LABEL maintainer=marjamis@amazon.com \
      version=testing
USER nobody
COPY --from=build .
ENTRYPOINT ["./app"]
