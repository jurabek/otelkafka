FROM golang:1.22-alpine AS builder

ENV PATH="/go/bin:${PATH}"
ENV GO111MODULE=on
ENV CGO_ENABLED=1
ENV GOPROXY=direct

RUN apk -U add ca-certificates
RUN apk update && apk upgrade && apk add pkgconf git bash build-base sudo librdkafka-dev

WORKDIR /go/src

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -tags musl --ldflags "-extldflags -static" -o consumer ./example/consumer

FROM scratch AS runner

COPY --from=builder /go/src/consumer /

EXPOSE 8080

ENTRYPOINT ["./consumer"]
