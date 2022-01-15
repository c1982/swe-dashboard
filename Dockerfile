FROM golang:1.17.0-buster AS builder

ADD ./ /go/src/swe-dashboard
WORKDIR /go/src/swe-dashboard
RUN go mod verify && \
    go mod vendor

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o swed ./app/swed

FROM alpine
COPY --from=builder /go/src/swe-dashboard/swed /swed
ENTRYPOINT [ "/swed" ]