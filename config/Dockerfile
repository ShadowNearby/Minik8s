FROM docker.io/library/ubuntu:22.04

RUN apt update && apt install -y \
    curl \
    git \
    gnupg \
    lsb-release \
    sudo \
    wget \
    iproute2 \
    iputils-ping

CMD ["/bin/bash","-c", "while true; do echo hello world; sleep 10; done"]
    