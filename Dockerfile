FROM golang:1.19

WORKDIR /webapp

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY src/ ./src/

# ENV GIN_MODE=release
ENV API_PORT=3000

RUN go build -o ./wee ./cmd/main.go

EXPOSE 3000

CMD [ "./wee" ]
