FROM golang:1.27.1-alpine3.24@sha256:8a5910f31396cd4d89662f56c68b3ae31d374308270a1c3bd96672ee5ed43414 AS builder
ADD . /go/lynis_exporter/
WORKDIR /go/lynis_exporter
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/lynis_exporter

FROM alpine:3.24.2@sha256:294b683cb724975bec92580e1e685676bd4b50bda910ddb8c51d4cabeaec77e6
RUN apk --no-cache add ca-certificates bash
COPY --from=builder /go/bin/lynis_exporter /app/lynis_exporter
EXPOSE 9730
ENV ARGS="-config.path=/config/config.yml"
VOLUME /report_data
VOLUME /config
ENTRYPOINT /app/lynis_exporter ${ARGS}
