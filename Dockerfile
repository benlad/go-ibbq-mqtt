FROM alpine:3.7

COPY go-ibbq-mqtt /app/go-ibbq-mqtt

CMD ["LOGXI=*=INF" "/app/go-ibbq-mqtt"]
