FROM golang:1.21-alpine AS builder
ENV GO111MODULE=on
RUN apk add git
WORKDIR /go/src/app
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY . .
RUN apk add --no-cache bash ca-certificates git gcc g++ libc-dev
RUN cd cmd/server/ && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o server .

FROM alpine:3.20
WORKDIR /go/src/app
RUN apk add --no-cache \
    tzdata \
    curl

RUN cp /usr/share/zoneinfo/Asia/Ho_Chi_Minh /etc/localtime
RUN echo "Asia/Ho_Chi_Minh" > /etc/timezone
COPY --from=builder /go/src/app/cmd/server/server /server
COPY --from=builder /go/src/app/.env /go/src/app/.env

EXPOSE 7000
CMD ["/server"]
