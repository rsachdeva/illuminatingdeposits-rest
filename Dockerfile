FROM golang:alpine AS build
WORKDIR /illuminatingdeposits
COPY go.mod .
RUN go mod download
COPY internal/ ./internal/
COPY cmd/ ./cmd/
WORKDIR /illuminatingdeposits/cmd/deltacli
RUN go build
WORKDIR /illuminatingdeposits/cmd/deltaapi
RUN go build

FROM alpine
RUN apk update
RUN apk add bash
WORKDIR /cmd
COPY --from=build /illuminatingdeposits/cmd/deltacli/deltacli .
COPY --from=build /illuminatingdeposits/cmd/deltaapi/deltaapi .