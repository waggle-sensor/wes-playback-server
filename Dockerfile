FROM golang:1.14 AS builder
WORKDIR /build
COPY go.mod main.go ./
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /build/playback-server

FROM alpine:latest
COPY --from=builder /build/playback-server /playback-server
EXPOSE 8090
ENTRYPOINT [ "/playback-server" ]
CMD [ "-data", "/data" ]
