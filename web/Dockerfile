FROM golang:1.14.3-alpine AS builder

RUN mkdir /app
ADD . /app
WORKDIR /app
RUN apk add --no-cache git \
                openssh-client \
                ca-certificates
RUN git config --global url."https://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:@github.com/".insteadOf "https://github.com/"
RUN export GO111MODULE=on
RUN export GOPRIVATE=github.com/cody6750/Web-Crawler
RUN go get github.com/cody6750/codywebapi/webServer/Handler
RUN go get github.com/cody6750/codywebapi/webCrawler
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./...

FROM alpine:latest AS production
COPY --from=builder /app .
CMD ["./main"]