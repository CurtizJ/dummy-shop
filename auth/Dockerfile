FROM golang:latest
WORKDIR /auth
COPY go.mod .
RUN go mod download
COPY . .
COPY . .

# Generate docs.
# RUN go get -u github.com/swaggo/swag/cmd/swag
# RUN swag init

RUN go build -o main .
CMD ["./main"]
