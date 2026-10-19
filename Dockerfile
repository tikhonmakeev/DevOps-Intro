FROM golang:1.25.13-alpine3.23 AS build
WORKDIR /src
COPY app/ .
RUN CGO_ENABLED=0 go build -trimpath -o /quicknotes .

FROM alpine:3.23.6
RUN addgroup -g 10001 quicknotes && adduser -D -H -u 10001 -G quicknotes quicknotes \
    && mkdir /data && chown quicknotes:quicknotes /data
COPY --from=build /quicknotes /usr/local/bin/quicknotes
COPY app/seed.json /seed.json
USER quicknotes
ENV ADDR=:8080 DATA_PATH=/data/notes.json SEED_PATH=/seed.json
EXPOSE 8080
HEALTHCHECK --interval=10s --timeout=3s --retries=3 CMD wget -q -O /dev/null http://127.0.0.1:8080/health || exit 1
ENTRYPOINT ["quicknotes"]
