FROM golang:1.19

WORKDIR /webapp

COPY go.mod go.sum ./
RUN go mod download

COPY app/ ./app/
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY views/ ./views/

# ENV GIN_MODE=release

RUN go build -o ./wee ./cmd/main.go

EXPOSE 3000

CMD [ "./wee" ]
