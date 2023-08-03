FROM alpine:latest
WORKDIR /app
COPY loggerService .
CMD ["/app/loggerService"]