FROM golang:1.19-alpine

WORKDIR /app

COPY src/* ./
RUN go mod download && \
    go build -o /app/directory-watcher

CMD [ "/app/directory-watcher" ]
