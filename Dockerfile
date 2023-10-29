FROM golang:1.21.3-alpine3.18 AS builder

RUN apk update && apk add --no-cache git

WORKDIR $GOPATH/src/smart-roast/backend/
COPY . .

RUN go mod tidy

# RUN go get -d -v

RUN go build -o /go/bin/smart-roast cmd/platform/main.go

FROM alpine:latest
RUN apk add --update tzdata
ENV TZ=Asia/Jakarta
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
RUN rm -rf /var/cache/apk/*
COPY --from=builder /go/bin/smart-roast /go/bin/smart-roast

EXPOSE 3000
EXPOSE 5432

ENTRYPOINT [ "/go/bin/smart-roast" ]
