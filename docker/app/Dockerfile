FROM golang:1.19

# RUN  go install github.com/cosmtrek/air@latest
RUN go install github.com/cespare/reflex@latest
# COPY docker/app/.air.toml /
COPY docker/app/reflex.conf /
COPY internal/ /internal


ENTRYPOINT ["reflex", "-c", "/reflex.conf"]
