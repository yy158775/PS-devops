#build stage
FROM golang:1.16 AS builder
# RUN set -eux && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
# RUN apk --no-cache add ca-certificates
RUN go env -w GOPROXY=https://goproxy.cn,direct
WORKDIR /
COPY code/helloworld .
RUN go mod download
RUN go build  -o /app -v ./greeter_server/main.go

#final stage
FROM scratch
# 看看上面这行能不能改掉
# RUN apk --no-cache add ca-certificates
COPY --from=builder /app /app
WORKDIR /
EXPOSE 50051
CMD ["./app"]

#FROM alpine:latest
#RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
#RUN apk --no-cache add ca-certificates
#COPY app /app
#WORKDIR /
#ENTRYPOINT /app