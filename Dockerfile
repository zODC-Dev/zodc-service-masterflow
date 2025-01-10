FROM golang:1.23.4-alpine
WORKDIR /app
COPY . .
RUN go mod download
EXPOSE 8080
CMD ["go", "run", "./src/cmd/main.go"]