# ============================== BINARY BUILDER ==============================
FROM golang:1.12.1 as builder

# Copy in the source
COPY . /src
WORKDIR /src

# Download and Install Tensorflow CPU
RUN mkdir local && \
    curl -L https://storage.googleapis.com/tensorflow/libtensorflow/libtensorflow-cpu-linux-x86_64-1.12.0.tar.gz | \
    tar -C local -xz && \
    cp -a local /usr

# Vendor, Test and Build the Binary
RUN GO111MODULE=on go mod vendor
RUN go test ./...
RUN CGO_ENABLED=1 go build -o ./bin/server

# ================================ FINAL IMAGE ================================
FROM ubuntu:18.04

# Dependencies
RUN apt-get update -y && \
    apt-get upgrade -y

# Copy Tensorflow
COPY --from=builder /src/local /usr
ENV LIBRARY_PATH $LIBRARY_PATH:/usr/local/lib
ENV LD_LIBRARY_PATH $LD_LIBRARY_PATH:/usr/local/lib

# Copy Graph/Labels, static files and Binary
COPY --from=builder /src/graph /graph
COPY --from=builder /src/service/templates /service/templates
COPY --from=builder /src/service/static /service/static
COPY --from=builder /src/bin/server /usr/local/bin/server
CMD ["server"]