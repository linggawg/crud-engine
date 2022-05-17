FROM golang:alpine

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY . .

RUN go mod tidy && go install && go build -o binary

CMD ["/app/binary"]