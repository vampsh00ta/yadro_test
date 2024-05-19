FROM golang:latest
WORKDIR /app
COPY  go.mod    ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build  -o task.exe
ENTRYPOINT  ["./task.exe"]
