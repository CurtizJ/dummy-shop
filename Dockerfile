FROM golang:latest
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .

# Generate docs.
RUN go get -u github.com/swaggo/swag/cmd/swag
RUN swag init

RUN go build -o main .
CMD ["./main"]
