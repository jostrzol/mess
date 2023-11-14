FROM golang:1.21

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY . .

# Build
ENV RELEASE=true
RUN CGO_ENABLED=0 GOOS=linux go build -o /mess-backend ./cmd/mess-server

EXPOSE 4000

# Run
CMD ["/mess-backend"]
