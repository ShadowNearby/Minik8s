FROM ubuntu:latest

COPY ./system_load /system_load

RUN apt update && apt install stress

EXPOSE 8080

ENTRYPOINT ["/system_load"]
