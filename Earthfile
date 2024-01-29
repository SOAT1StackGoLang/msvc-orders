VERSION 0.7
all:
    BUILD --platform=linux/amd64 --platform=linux/arm64 +msvc
    BUILD --platform=linux/amd64 --platform=linux/arm64 +debug
amd64:
    BUILD --platform=linux/amd64 +msvc
    BUILD --platform=linux/amd64 +debug
arm64:
    BUILD --platform=linux/arm64 +msvc
    BUILD --platform=linux/arm64 +debug
file:
    LOCALLY
    SAVE ARTIFACT ./
deps:
    FROM golang:alpine
    WORKDIR /build
    COPY +file/* ./
    RUN ls -althR
    RUN apk add --no-cache git
    RUN go mod tidy
    RUN go mod download
    #RUN go get -u github.com/swaggo/swag/cmd/swag
    #RUN go install github.com/swaggo/swag/cmd/swag
    #RUN swag init -g ../../cmd/web/routes.go -o ./docs -d ./internal/handlers

compile:
    FROM +deps
    ARG GOOS=linux
    ARG GOARCH=amd64
    ARG VARIANT
    RUN ls -alth && pwd
    RUN GOARM=${VARIANT#v} CGO_ENABLED=0 go build \
        -installsuffix 'static' \
        -o compile/app cmd/server/*.go
    RUN GOARM=${VARIANT#v} CGO_ENABLED=0 go build \
        -installsuffix 'static' \
        -o compile/migs cmd/migrations/*.go
    SAVE ARTIFACT compile/app /app AS LOCAL compile/app
    SAVE ARTIFACT compile/migs /migs AS LOCAL compile/migs
#--ldflags "-X 'msvc.Version=v0.0.3' -X 'msvc.BuildTime=$(date "+%H:%M:%S--%d/%m/%Y")' -X 'msvc.GitCommit=$(git rev-parse --short HEAD)'" \

migrations:
    ARG EARTHLY_TARGET_TAG_DOCKER
    ARG EARTHLY_GIT_SHORT_HASH
    ARG TARGETPLATFORM
    ARG TARGETARCH
    ARG TARGETVARIANT
    FROM --platform=$TARGETPLATFORM gcr.io/distroless/static
    ## enable multiple debug version with shell
    #FROM --platform=$TARGETPLATFORM gcr.io/distroless/static:debug
    #FROM --platform=$TARGETPLATFORM alpine:latest
    LABEL org.opencontainers.image.source=https://github.com/soat1stackgolang/tech-challenge
    LABEL org.opencontainers.image.description="Migrations Image only have the migrations binary and nothing else"
    WORKDIR /
    COPY \
        --platform=linux/amd64 \
        (+compile/migs --GOARCH=$TARGETARCH --VARIANT=$TARGETVARIANT) /migs
    ENV GIN_MODE=release
    ENTRYPOINT ["/migs"]
    EXPOSE 8080
    SAVE IMAGE --push ghcr.io/soat1stackgolang/msvc-orders:migs-$EARTHLY_TARGET_TAG_DOCKER
    SAVE IMAGE --push ghcr.io/soat1stackgolang/msvc-orders:migs-$EARTHLY_GIT_SHORT_HASH


msvc:
    ARG EARTHLY_TARGET_TAG_DOCKER
    ARG EARTHLY_GIT_SHORT_HASH
    ARG TARGETPLATFORM
    ARG TARGETARCH
    ARG TARGETVARIANT
    FROM --platform=$TARGETPLATFORM gcr.io/distroless/static
    ## enable multiple debug version with shell
    #FROM --platform=$TARGETPLATFORM gcr.io/distroless/static:debug
    #FROM --platform=$TARGETPLATFORM alpine:latest
    LABEL org.opencontainers.image.source=https://github.com/soat1stackgolang/msvc-orders
    LABEL org.opencontainers.image.description="Main App Image only have the app binary and nothing else"
    WORKDIR /
    COPY \
        --platform=linux/amd64 \
        (+compile/app --GOARCH=$TARGETARCH --VARIANT=$TARGETVARIANT) /app
    ENV GIN_MODE=release
    ENTRYPOINT ["/app"]
    EXPOSE 8080
    SAVE IMAGE --push ghcr.io/soat1stackgolang/msvc-orders:msvc-$EARTHLY_TARGET_TAG_DOCKER
    SAVE IMAGE --push ghcr.io/soat1stackgolang/msvc-orders:msvc-$EARTHLY_GIT_SHORT_HASH

debug:
    ARG EARTHLY_TARGET_TAG_DOCKER
    ARG EARTHLY_GIT_SHORT_HASH
    ARG TARGETPLATFORM
    ARG TARGETARCH
    ARG TARGETVARIANT
    #FROM --platform=$TARGETPLATFORM gcr.io/distroless/static
    ## enable multiple debug version with shell
    #FROM --platform=$TARGETPLATFORM gcr.io/distroless/static:debug
    FROM --platform=$TARGETPLATFORM alpine:latest
    LABEL org.opencontainers.image.source=https://github.com/soat1stackgolang/msvc-orders
    LABEL org.opencontainers.image.description="Debug Image will have all binaries"
    WORKDIR /
    COPY \
        --platform=linux/amd64 \
        (+compile/app --GOARCH=$TARGETARCH --VARIANT=$TARGETVARIANT) /app
    COPY \
        --platform=linux/amd64 \
        (+compile/migs --GOARCH=$TARGETARCH --VARIANT=$TARGETVARIANT) /migs
    ENV GIN_MODE=release
    CMD /migs; /app
    EXPOSE 8080
    SAVE IMAGE --push ghcr.io/soat1stackgolang/msvc-orders:debug-$EARTHLY_TARGET_TAG_DOCKER
    SAVE IMAGE --push ghcr.io/soat1stackgolang/msvc-orders:debug-$EARTHLY_GIT_SHORT_HASH
