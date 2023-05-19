FROM alpine:latest

RUN set -eux ; \
    apk add ca-certificates tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime ; \
    echo "Asia/Shanghai" > /etc/timezone ; \
    apk del tzdata

WORKDIR /app

COPY templates  /app/templates
COPY README.md /app/
COPY webtool /app/

EXPOSE 20231

CMD ["./20231"]
