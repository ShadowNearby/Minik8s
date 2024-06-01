FROM python:3.9-slim

WORKDIR /app

RUN pip install flask

COPY app.py /app/app.py

EXPOSE 80

CMD ["python", "app.py"]
