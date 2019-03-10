# ============================== BINARY BUILDER ==============================
FROM golang:1.12 as builder

# Copy in the source
WORKDIR /go/src/s32x.com/tfclass
COPY / .

# Dependencies
RUN apt-get update -y && \
    apt-get upgrade -y
RUN wget https://storage.googleapis.com/tensorflow/libtensorflow/libtensorflow-cpu-linux-x86_64-1.12.0.tar.gz
RUN tar -C /usr/local -xzf libtensorflow-cpu-linux-x86_64-1.12.0.tar.gz
RUN rm -f libtensorflow-cpu-linux-x86_64-1.12.0.tar.gz
RUN export GO111MODULE=on && \
    go get github.com/ahmetb/govvv && \
    go mod vendor

# Build the binary
RUN CGO_ENABLED=1 GOOS=linux go build -o ./bin/server

# ================================ FINAL IMAGE ================================
FROM ubuntu:latest

# Dependencies
RUN apt-get update -y && \
    apt-get upgrade -y && \
    apt-get install -y wget
RUN wget https://storage.googleapis.com/tensorflow/libtensorflow/libtensorflow-cpu-linux-x86_64-1.12.0.tar.gz
RUN tar -C /usr/local -xzf libtensorflow-cpu-linux-x86_64-1.12.0.tar.gz
RUN rm -f libtensorflow-cpu-linux-x86_64-1.12.0.tar.gz

# Graph/Labels files
COPY graph /graph
COPY test /test
COPY service/static /service/static

# Environment
ENV LIBRARY_PATH=$LIBRARY_PATH:/usr/local/lib
ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/usr/local/lib

# Binary
COPY --from=builder /go/src/s32x.com/tfclass/bin/server /usr/local/bin/server
CMD ["server"]