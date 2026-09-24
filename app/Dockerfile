FROM golang:1.24-alpine AS builder
RUN mkdir /data
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /quicknotes

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /quicknotes /quicknotes
COPY seed.json /seed.json
COPY --from=builder --chown=nonroot:nonroot /data /data
ENV ADDR=:8080 DATA_PATH=/data/notes.json SEED_PATH=/seed.json
EXPOSE 8080
USER nonroot
ENTRYPOINT ["/quicknotes"]
