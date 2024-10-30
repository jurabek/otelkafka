FROM golang:1.22-alpine AS builder

ENV PATH="/go/bin:${PATH}"
ENV GO111MODULE=on
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOPROXY=direct

WORKDIR /go/src

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN apk -U add ca-certificates
RUN apk update && apk upgrade && apk add pkgconf git bash build-base sudo
RUN apk update && apk add librdkafka-dev

RUN go build -tags musl --ldflags "-extldflags -static" -o producer ./example/producer

FROM scratch AS runner

COPY --from=builder /go/src/producer /

EXPOSE 8080

ENTRYPOINT ["./producer"]