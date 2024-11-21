FROM golang:1.22-alpine AS builder

WORKDIR /usr/local/src

RUN apk --no-cache add bash git make gcc gettext musl-dev

# dependencies
COPY go.mod go.sum ./
RUN go mod download

# build
COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -tags netgo -tags musl -o ./bin/app cmd/tinder-geo/main.go

FROM alpine AS runner

COPY --from=builder /usr/local/src/bin/app /
COPY config/config.yaml /config.yaml

ENV CONFIG_PATH=config.yaml

CMD ["/app"]