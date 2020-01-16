FROM golang:rc-alpine AS build-env
ADD . /app
WORKDIR /app
# install dependencies
RUN apk add --no-cache git \
    git \
    binutils-gold \
    curl \
    g++ \
    gcc \
    gnupg \
    libgcc \
    linux-headers \
    make
# build
RUN make build

#FROM golang:rc-alpine
#COPY --from=build-env /app/build/aws-es-proxy-go /usr/local/bin/aws-es-proxy-go
#ENTRYPOINT ["/usr/local/bin/aws-es-proxy-go"]
