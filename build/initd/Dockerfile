ARG SEMVER="undefined-docker"
ARG COMMITSHA="undefined-docker"
ARG BUILDTIME="undefined-docker"

FROM golang:1.26.1-alpine AS init-build
ARG INIT_BUILDFLAGS="-v"
ARG INIT_LDFLAGS="-s -w"
ARG SEMVER
ARG COMMITSHA
ARG BUILDTIME

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    VERSIONPKG="github.com/stevezaluk-org/apathy-os/internal/version"

WORKDIR /build

# Needs docker-ignore
COPY . .

RUN go build -o init \
    $INIT_BUILDFLAGS \
    -ldflags="$INIT_LDFLAGS -X $VERSIONPKG.SemVer=$SEMVER -X $VERSIONPKG.CommitSHA=$COMMITSHA -X $VERSIONPKG.BuildDate=$BUILDTIME" \
    ./build/initd && \
    chmod +x init

FROM gcr.io/distroless/static-debian13:latest

WORKDIR /app
COPY --from=init-build /build/init /app/init

ENTRYPOINT ["init"]
