FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY app/ ./
COPY internal/ ./
COPY views/ ./

RUN go build -o /wee app/main.go

EXPOSE 3000

CMD [ "/wee" ]
