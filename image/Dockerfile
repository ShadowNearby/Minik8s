FROM pytorch/pytorch:latest

COPY ./  /app

WORKDIR /app
# 换成中国的源
RUN pip3 config set global.index-url https://pypi.tuna.tsinghua.edu.cn/simple/ 

RUN pip3 install -r ./requirements.txt

EXPOSE 18080

ENTRYPOINT ["python3"]

CMD ["server.py"]
