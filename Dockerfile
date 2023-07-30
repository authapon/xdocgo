FROM golang:alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o xdocgo .

FROM scratch

COPY --from=builder ["/build/xdocgo", "/"]
ADD ./vue/dist /vue/dist
ADD ./template /template
ADD xdoc.yml /etc/xdoc.yml

EXPOSE 80

ENTRYPOINT ["/xdocgo"]
