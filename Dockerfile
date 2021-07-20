FROM golang as builder
ADD . /go/lynis_exporter/
WORKDIR /go/lynis_exporter
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/lynis_exporter

FROM alpine:latest
RUN apk --no-cache add ca-certificates bash
COPY --from=builder /go/bin/lynis_exporter /app/lynis_exporter
EXPOSE 9730
ENV ARGS "-config.path=/config/config.yml"
VOLUME /report_data
VOLUME /config
ENTRYPOINT /app/lynis_exporter ${ARGS}