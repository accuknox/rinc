FROM golang:1.23.3-alpine3.20 as deps
RUN apk add --update --no-cache ca-certificates git
ENV GOPATH=/go
WORKDIR /deps
COPY go.mod /deps
COPY go.sum /deps
RUN go mod download

FROM golang:1.23.3-alpine3.20 as builder
COPY --from=deps /go /go
ENV GOPATH=/go
ENV CGO_ENABLED=0
COPY . /rinc
WORKDIR /rinc
RUN mkdir bin -p && go build -o bin/main ./cmd/rinc

FROM alpine:3.20 as runner
COPY --from=builder /rinc/bin/main /rinc/main
COPY static /rinc/static
WORKDIR /rinc
RUN adduser --disabled-password --no-create-home rinc
RUN mkdir -p reports && chown rinc:rinc reports/
USER rinc
ENTRYPOINT [ "./main" ]
