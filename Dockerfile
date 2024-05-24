FROM golang:1.20-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOARCH=amd64 GOOS=linux go build -o ./analyzerbin ./cmd

FROM alpine:3.20.0
WORKDIR /app
COPY --from=build /app/analyzerbin  ./analyzerbin
COPY templates  ./templates/
COPY config  ./config/

ENTRYPOINT ["./analyzerbin"]
