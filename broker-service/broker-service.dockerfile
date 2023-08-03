FROM alpine:latest
WORKDIR /app
COPY brokerService .
CMD ["/app/brokerService"]