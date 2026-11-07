ARG SEMVER="undefined-docker"
ARG COMMITSHA="undefined-docker"
ARG BUILDTIME="undefined-docker"

FROM golang:1.26.1-alpine AS serviced-build
ARG SERVICED_BUILDFLAGS="-v"
ARG SERVICED_LDFLAGS="-s -w"
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

RUN go build -o serviced \
    $SERVICED_BUILDFLAGS \
    -ldflags="$SERVICED_LDFLAGS -X $VERSIONPKG.SemVer=$SEMVER -X $VERSIONPKG.CommitSHA=$COMMITSHA -X $VERSIONPKG.BuildDate=$BUILDTIME" \
    ./build/serviced && \
    chmod +x serviced

FROM gcr.io/distroless/static-debian13:latest

WORKDIR /app
COPY --from=serviced-build /build/serviced /app/serviced

ENTRYPOINT ["/app/serviced"]
