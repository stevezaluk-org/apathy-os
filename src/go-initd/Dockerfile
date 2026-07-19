FROM golang:1.26.1-alpine AS init-build
ARG INIT_BUILDFLAGS="-v"
ARG INIT_LDFLAGS="-s -w"

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY ./go-initd go-initd/
COPY ./go-manifest go-manifest/

# this can be removed once this goes public
RUN go work init ./go-initd ./go-manifest

RUN go build -o init $INIT_BUILDFLAGS -ldflags="$INIT_LDFLAGS" ./go-initd && chmod +x init

FROM gcr.io/distroless/static-debian13:latest

WORKDIR /app
COPY --from=init-build /build/init /app/init

ENTRYPOINT ["init"]