FROM golang:alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY *.go ./
COPY ./templates/ ./templates/
# COPY ./config ./config

RUN go build -o /ckb

EXPOSE 8081

CMD [ "/ckb" ]
