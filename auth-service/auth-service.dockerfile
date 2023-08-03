FROM alpine:latest
WORKDIR /app
COPY authService .
CMD ["/app/authService"]