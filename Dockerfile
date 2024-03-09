ARG GO_VERSION=1.22
FROM golang:${GO_VERSION}-alpine as builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /baleen.app ./cmd/baleen.app


FROM alpine:latest
COPY --from=builder /baleen.app /usr/local/bin/
CMD ["baleen.app", "-migrate-db"]
