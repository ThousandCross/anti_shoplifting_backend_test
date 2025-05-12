FROM golang:1.18


WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN apt-get update \
    && apt-get install -y vim \
    && go install github.com/cosmtrek/air@latest

# CMD ["go", "run", "main.go"]
# Live reload for Go apps
CMD ["air"]
