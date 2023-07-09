FROM golang:1.20-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY *.go ./
RUN go install

COPY *.go ./

RUN go build -o /nanodns

EXPOSE 53

CMD [ "/nanodns" ]