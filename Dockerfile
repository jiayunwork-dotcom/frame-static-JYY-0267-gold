FROM golang:1.21-alpine
ENV GOTOOLCHAIN=local
ENV CGO_ENABLED=0
ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build ./...
CMD ["/bin/sh"]
