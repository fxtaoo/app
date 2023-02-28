FROM alpine:latest

RUN set -eux; \
    apk add ca-certificates tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata

WORKDIR /app

COPY templates  /app/templates
COPY README.md /app/
COPY main /app/

EXPOSE 54328

CMD ["./main"]
