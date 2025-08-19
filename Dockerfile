FROM golang:alpine AS build
WORKDIR /go/src/app
COPY . .
RUN export GO111MODULE=on && export GOPROXY=https://goproxy.cn && \
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ipproxypool -a -ldflags "-s -w" -tags timetzdata main.go && \
cp ipproxypool /


FROM alpine
COPY --from=build /ipproxypool /
ENTRYPOINT ["/ipproxypool"]
EXPOSE 6060
