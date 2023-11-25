FROM golang

WORKDIR /app

ENV TZ Asia/Tehran

COPY ../go.mod ../go.sum ./
RUN go mod download

COPY ../. .
RUN go build -o main cmd/main.go

EXPOSE 5000

CMD ["./main"]