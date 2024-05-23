FROM python:3.9-alpine

RUN pip install prometheus-client

COPY prometheus_example_program.py /program.py

CMD ["python", "/program.py"]