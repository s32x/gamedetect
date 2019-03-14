# ============================== BINARY BUILDER ==============================
FROM golang:latest as builder

# Copy in the source
COPY . /service
WORKDIR /service

# Dependencies
RUN GO111MODULE=on go mod vendor

# Install Tensorflow
RUN wget https://storage.googleapis.com/tensorflow/libtensorflow/libtensorflow-cpu-linux-x86_64-1.12.0.tar.gz
RUN tar -C /usr/local -xzf libtensorflow-cpu-linux-x86_64-1.12.0.tar.gz

# Test before building
RUN go test ./...

# Build the binary
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o ./bin/server

# ================================ FINAL IMAGE ================================
FROM ubuntu:latest

# Dependencies
RUN apt-get update -y && \
    apt-get upgrade -y && \
    apt-get install -y wget

# Install Tensorflow
RUN wget https://storage.googleapis.com/tensorflow/libtensorflow/libtensorflow-cpu-linux-x86_64-1.12.0.tar.gz
RUN tar -C /usr/local -xzf libtensorflow-cpu-linux-x86_64-1.12.0.tar.gz
RUN rm -f libtensorflow-cpu-linux-x86_64-1.12.0.tar.gz

# Graph/Labels and static files
COPY graph /graph
COPY service/templates /service/templates
COPY service/static /service/static

# Environment
ENV LIBRARY_PATH=$LIBRARY_PATH:/usr/local/lib
ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/usr/local/lib

# Binary
COPY --from=builder /service/bin/server /usr/local/bin/server
CMD ["server"]