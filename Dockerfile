FROM --platform=${BUILDPLATFORM} golang:1.21-alpine AS builder
LABEL maintainer="Tom Helander <thomas.helander@gmail.com>"

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS TARGETARCH
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o purpleair_exporter .

FROM alpine:3.18.4
LABEL maintainer="Tom Helander <thomas.helander@gmail.com>"

WORKDIR /app

COPY --from=builder /src/purpleair_exporter .

EXPOSE 9811

ENTRYPOINT ["/app/purpleair_exporter"]
