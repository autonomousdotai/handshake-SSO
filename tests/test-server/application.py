from flask import Flask, request
import os
import sys
import json
import tempfile

app = Flask(__name__)

@app.route("/", methods=["GET", "POST"])
def hello():
    req_data = {} 
    req_data['endpoint'] = request.endpoint
    req_data['method'] = request.method
    req_data['cookies'] = request.cookies
    req_data['data'] = request.data
    req_data['headers'] = dict(request.headers)
    req_data['headers'].pop('Cookie', None)
    req_data['args'] = request.args
    req_data['form'] = request.form
    req_data['json'] = request.json
    req_data['remote_addr'] = request.remote_addr
    
    files = []
    for name, fs in request.files.iteritems():
        dst = tempfile.NamedTemporaryFile()
        fs.save(dst)
        dst.flush()
        filesize = os.stat(dst.name).st_size
        dst.close()
        files.append({'name': name, 'filename': fs.filename, 'filesize': filesize, 'mimetype': fs.mimetype, 'mimetype_params': fs.mimetype_params})

    req_data['files'] = files
    print req_data
    return "Hello world!"
