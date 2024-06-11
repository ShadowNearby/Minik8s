FROM python:3.9-slim

RUN apt-get update
RUN apt-get -y install openssh-server python3-pip vim

WORKDIR /root
RUN pip3 install  paramiko scp


COPY ./id_rsa .ssh/id_rsa

COPY ./known_hosts .ssh/known_hosts

COPY ./job.py job.py




