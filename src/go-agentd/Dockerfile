FROM golang:1.26.1-alpine AS agentd-build
ARG AGENT_BUILDFLAGS="-v"
ARG AGENT_LDFLAGS="-s -w"

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY ./go-agentd agentd/
COPY ./go-manifest go-manifest/

# this can be removed once this goes public
RUN go work init ./agentd ./go-manifest

RUN go build -o agent $AGENT_BUILDFLAGS -ldflags="$AGENT_LDFLAGS" ./agentd && chmod +x agent

FROM gcr.io/distroless/static-debian13:latest

WORKDIR /app
COPY --from=agentd-build /build/agent /app/agentd

ENTRYPOINT ["/app/agentd"]