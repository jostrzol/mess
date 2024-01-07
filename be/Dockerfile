FROM golang:1.21

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV RELEASE=true
RUN CGO_ENABLED=0 GOOS=linux go build -o /mess-backend ./cmd/mess-server

EXPOSE 4000

CMD ["/mess-backend"]
