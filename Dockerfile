FROM ubuntu:xenial AS lotus-builder

RUN apt update -y
RUN apt install -y mesa-opencl-icd ocl-icd-opencl-dev gcc git bzr jq pkg-config curl wget nano
RUN apt upgrade -y

RUN wget -c https://golang.org/dl/go1.14.6.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.14.6.linux-amd64.tar.gz
ENV PATH=$PATH:/usr/local/go/bin
ENV LOTUS_SKIP_GENESIS_CHECK=_yes_
WORKDIR /app/lotus

RUN git clone https://github.com/filecoin-project/lotus .
RUN git checkout v0.5.4
RUN make 2k

# Cache the params (~1 GB download)
RUN ./lotus fetch-params 2048

WORKDIR /app/go-fil-markets
COPY . /app/go-fil-markets

# Lotus API port
EXPOSE 1234
# Miner API port
EXPOSE 2345
# Market API port
EXPOSE 8888

CMD ["./run_docker.sh"]