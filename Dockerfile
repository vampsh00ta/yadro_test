FROM golang:latest
WORKDIR /app
COPY  go.mod  *.txt  ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build  -o task.exe
CMD ["./task.exe", "file.txt"]
