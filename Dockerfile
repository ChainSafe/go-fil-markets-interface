FROM ubuntu:xenial AS lotus-builder

RUN apt update -y
RUN apt install -y mesa-opencl-icd ocl-icd-opencl-dev gcc git bzr jq pkg-config curl wget nano libhwloc-dev
RUN apt upgrade -y

RUN wget -c https://golang.org/dl/go1.16.3.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.16.3.linux-amd64.tar.gz
ENV PATH=$PATH:/usr/local/go/bin
ENV LOTUS_SKIP_GENESIS_CHECK=_yes_
ENV DOCKER=_yes_
ENV CGO_CFLAGS_ALLOW="-D__BLST_PORTABLE__"
ENV CGO_CFLAGS="-D__BLST_PORTABLE__"
ENV LOTUS_PATH=$HOME/.lotusDevnet
ENV LOTUS_MINER_PATH=$HOME/.lotusminerDevnet
WORKDIR /app/lotus

RUN git clone https://github.com/filecoin-project/lotus .
RUN git checkout v1.9.0
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