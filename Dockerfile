FROM golang:1.10

WORKDIR /go/app

RUN go get -u github.com/ikasamah/homecast
RUN go get -u github.com/joho/godotenv


