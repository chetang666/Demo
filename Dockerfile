FROM golang:latest

RUN mkdir /app

ADD . /app

WORKDIR /app

# Fetch dependencies.
RUN go mod download

# we run go build to compile the binary
# executable of our Go program
RUN go build -o main .

# our newly created binary executable
CMD ["/app/main"]