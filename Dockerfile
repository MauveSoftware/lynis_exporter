FROM golang:1.27.2-alpine3.24@sha256:85dc1069ac644ea3c527b177303a406eb3358192816cd7f9e5848eb658851673 AS builder
ADD . /go/lynis_exporter/
WORKDIR /go/lynis_exporter
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/lynis_exporter

FROM alpine:3.24.3@sha256:be2656061b536c1883c4034790e21d0a09c5aaaf51968b89368a3c0371ea2c9a
RUN apk --no-cache add ca-certificates bash
COPY --from=builder /go/bin/lynis_exporter /app/lynis_exporter
EXPOSE 9730
ENV ARGS="-config.path=/config/config.yml"
VOLUME /report_data
VOLUME /config
ENTRYPOINT /app/lynis_exporter ${ARGS}
