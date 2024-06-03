import os
import func
import json

from flask import Flask, request

app = Flask(__name__)
host_conf = "0.0.0.0"
port_conf = os.environ.get('PORT', 18080)


@app.route('/', methods=['GET'])
def welecome():
    welecomeWords = "Welcome to use this server!"
    usage = "Usage: send put/post request to this url!"
    return welecomeWords + '\n' + usage

@app.route("/", methods=['POST'])
def callCloudFuncByPost():
    params = request.json
    headers = {'Content-Type': 'text/plain'}
    try:
        result = func.run(**params)
        response = Response(str(result), headers=headers, status=200)
        return response
    except TypeError as e:
        response = Response("", headers=headers, status=200)
        return response
    except Exception as e:
        response = Response(str(e), headers=headers, status=500)
        return str(e)


@app.route("/config", methods=['GET'])
def getConfig():
    config = {
        "host": host_conf,
        "port": port_conf,
    }
    return json.dumps(config)

if __name__ == '__main__':
    app.run(host=host_conf, port=port_conf, debug=False)