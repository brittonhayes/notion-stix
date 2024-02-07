ARG GO_VERSION=1.21

# STAGE 1
FROM golang:${GO_VERSION}-alpine AS builder

WORKDIR /src/

COPY go.mod ./
RUN go mod download

COPY . /src/
RUN CGO_ENABLED=0 go build -o /bin/stix cmd/stix/main.go

# STAGE 2
FROM gcr.io/distroless/static-debian11:nonroot

LABEL maintainer="brittonhayes"
LABEL org.opencontainers.image.source="https://github.com/brittonhayes/notion-stix"
LABEL org.opencontainers.image.description="This is the Notion STIX Integration API. It allows you to integrate  threat intelligence data into Notion."
LABEL org.opencontainers.image.licenses="MIT"

COPY --from=builder --chown=nonroot:nonroot /bin/stix /bin/stix

EXPOSE 8080

ENTRYPOINT [ "/bin/stix" ]

CMD ["/bin/stix"]