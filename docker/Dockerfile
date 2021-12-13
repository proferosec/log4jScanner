FROM nginx:latest

RUN mkdir /certs

COPY nginx.conf /etc/nginx/nginx.conf

RUN openssl req -x509 -nodes -days 365 \
    -subj  "/CN=localhost" \
     -newkey rsa:2048 -keyout /certs/key.pem \
     -out /certs/cert.pem

