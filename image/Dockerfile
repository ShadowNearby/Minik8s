# Use the official Python image from the Docker Hub
FROM python:3.9-slim
LABEL authors="xinyun"

ENTRYPOINT ["top", "-b"]
# use basic image
# set workdir
WORKDIR /app
COPY . .

RUN apt-get update && \
    apt-get install -y python3-venv && \
    python3 -m venv /app/venv


EXPOSE 8081


# copy file into workdir
copy . .
# install dependencies from requirements.txt
RUN pip3 install -i https://pypi.tuna.tsinghua.edu.cn/simple --default-timeout=60 --no-cache-dir -r ./requirements.txt

# set env
ENV PATH="/usr/local/bin:${PATH}"

# run server.py, start the function server
CMD ["python3", "server.py"]