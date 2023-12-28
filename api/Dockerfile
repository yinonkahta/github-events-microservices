FROM golang:1.21-alpine

WORKDIR /app

COPY . .
WORKDIR api
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /api

EXPOSE 8080

CMD [ "/api" ]