FROM golang:1.18-bullseye

RUN go install github.com/cespare/reflex@v0.3.1
COPY ./docker/app/reflex.conf /

COPY ./docker/app/start.sh /

COPY ./internal /internal/

ENTRYPOINT ["reflex", "-c", "/reflex.conf"]
