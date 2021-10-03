FROM golang:1.16  AS builder
# RUN set -eux && sed -i 's/dl- cdn.alpinelinux.org/mirrors.ustc.edu.cn/g'
# RUN apk update && apk add --no-cache git
RUN go env -w GOPROXY=https://goproxy.cn,direct
WORKDIR /
COPY . .
RUN go mod download
# RUN go build  -o /app -v ./greeter_server/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app -ldflags="-w -s" ./greeter_server/main.go

FROM scratch
COPY --from=builder /app /app
EXPOSE 50051
ENTRYPOINT ["/app"]