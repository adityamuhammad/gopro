FROM golang:1.22.4-alpine

WORKDIR /app

COPY . .

RUN go build -o gopro/web cmd/web/main.go
RUN go build -o gopro/worker cmd/worker/main.go

RUN apk add --no-cache supervisor  # Install supervisord


# Copy supervisord.conf to /etc
COPY ./supervisord.conf /etc/supervisord.conf

EXPOSE 8080

RUN ls -l /app/gopro

CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]