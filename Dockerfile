FROM --platform=${BUILDPLATFORM} golang:1.21-alpine AS builder
LABEL maintainer="Tom Helander <thomas.helander@gmail.com>"

WORKDIR /app
COPY . .

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build .

FROM --platform=${BUILDPLATFORM} alpine:3.18.4
LABEL maintainer="Tom Helander <thomas.helander@gmail.com>"

WORKDIR /app

COPY --from=builder /app/purpleair_exporter .

EXPOSE 9811

ENTRYPOINT ["/app/purpleair_exporter"]
