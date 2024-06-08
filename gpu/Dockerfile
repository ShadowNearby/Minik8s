FROM ubuntu

RUN apt-get update
RUN apt-get -y install openssh-server python3-pip vim

RUN pip3 config set global.index-url https://pypi.tuna.tsinghua.edu.cn/simple/

RUN pip3 install -r ./requirements.txt

WORKDIR /root

RUN mkdir .ssh

COPY ./id_rsa .ssh/id_rsa

COPY ./known_hosts .ssh/known_hosts

COPY job.py addjob.py
