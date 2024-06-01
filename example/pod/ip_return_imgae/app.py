from flask import Flask
import os
import socket

app = Flask(__name__)

def get_ip():
    hostname = socket.gethostname()
    ip_address = socket.gethostbyname(hostname)
    return ip_address

@app.route('/')
def index():
    ip_address = get_ip()
    return f"Container IP Address: {ip_address}"

if __name__ == "__main__":
    app.run(host='0.0.0.0', port=80)
