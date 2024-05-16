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

CMD ["/bin/bash","-c", "tail -f /dev/null"]
    