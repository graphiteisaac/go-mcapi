FROM golang:1.20-alpine AS builder

ENV APP_HOME /go/src/storm
ENV CGO_ENABLED 0

WORKDIR /app
COPY go.mod .
COPY go.sum .

RUN go mod download
RUN go mod verify

COPY . .

RUN GOOS=linx GOARCH=amd64 go build -ldflags="-w -s" -o app .

FROM  alpine
COPY --from=builder $APP_HOME/app .

EXPOSE 3333
ENTRYPOINT ["./app"]

