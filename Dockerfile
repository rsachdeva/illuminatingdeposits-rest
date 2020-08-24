FROM golang:alpine AS build
WORKDIR /illuminatingdeposits
COPY go.mod .
RUN go mod download
COPY internal/ ./internal/
COPY cmd/ ./cmd/
WORKDIR /illuminatingdeposits/cmd/dbclient
RUN go build
WORKDIR /illuminatingdeposits/cmd/interestapi
RUN go build

FROM alpine
RUN apk update
RUN apk add bash
WORKDIR /cmd
COPY --from=build /illuminatingdeposits/cmd/dbclient/dbclient .
COPY --from=build /illuminatingdeposits/cmd/interestapi/interestapi .