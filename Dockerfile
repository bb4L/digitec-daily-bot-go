FROM golang:1.19 as builder

# Env variables
ENV GOOS linux
ENV CGO_ENABLED 0

WORKDIR /digitec_daily_bot
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o digitec-daily-bot

FROM alpine:3.16 as production

COPY --from=builder digitec_daily_bot/digitec-daily-bot .

CMD ./digitec-daily-bot