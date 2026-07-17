# syntax=docker/dockerfile:1.7
#
# Builds the single axisml-core image (Cluster Manager + Compute Service +
# Artifact Hub modules + the in-process Standalone runtime). The build
# context is this repository and the sibling AxisML checkout is supplied as the
# named `axisml` build context so local replace directives resolve:
#   docker build --build-context axisml=../axisml \
#     -f axisml-core/Dockerfile -t axisml-core:dev .

FROM --platform=$BUILDPLATFORM golang:1.26 AS builder
ARG TARGETOS
ARG TARGETARCH
WORKDIR /workspace

# Recreate the two sibling repository paths under /workspace. The main build
# context contains axisml-lite; the named context contains shared AxisML modules.
COPY --from=axisml pkg/axismlconfig/ axisml/pkg/axismlconfig/
COPY --from=axisml axisml-system/ axisml/axisml-system/
COPY axisml-core/ axisml-lite/axisml-core/

WORKDIR /workspace/axisml-lite/axisml-core
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -trimpath -ldflags='-s -w' \
        -o /out/axisml-core ./cmd/axisml-core

# axisml-core needs root to reach the Docker socket (design §7.2: it is the only
# process with socket access). static-debian12 (root) ships CA certs for OCI/S3.
FROM gcr.io/distroless/static-debian12:latest
COPY --from=builder /out/axisml-core /usr/local/bin/axisml-core
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/axisml-core"]
CMD ["serve"]
