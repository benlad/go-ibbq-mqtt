FROM arm32v7/golang:alpine as builder
RUN mkdir /build 
ADD . /build/
WORKDIR /build 
RUN GOOS=linux go build
FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/go-ibbq-mqtt /app/
WORKDIR /app
CMD ["/app/go-ibbq-mqtt"]
