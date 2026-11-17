FROM golang:1.27.4-alpine3.24@sha256:0a5f7fe796599ff487605987432e7619ae45beeb2fba4c2ce0b907879932aebf AS builder
ADD . /go/lynis_exporter/
WORKDIR /go/lynis_exporter
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/lynis_exporter

FROM alpine:3.24.4@sha256:f59abe656b692d901e184be3445e084786602108599d2a51d229c0f5a48c71bf
RUN apk --no-cache add ca-certificates bash
COPY --from=builder /go/bin/lynis_exporter /app/lynis_exporter
EXPOSE 9730
ENV ARGS="-config.path=/config/config.yml"
VOLUME /report_data
VOLUME /config
ENTRYPOINT /app/lynis_exporter ${ARGS}
