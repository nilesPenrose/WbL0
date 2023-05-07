# main.py

import requests
from flask import render_template, request, json

from app import app


@app.route('/', methods=['GET', 'POST'])
def index():
    print(request.form)
    if request.form.get("orders") and request.form["send"] == "send":
        r = requests.post(url="http://127.0.0.1:5000/order/get",
                          params={'orderUID': request.form.get("orders")})
        pj = str(json.dumps(r.json(), indent=2)).split('\n')

        return render_template('results.html', orders=pj)
    else:
        r = requests.get(url="http://127.0.0.1:5000/order/uid/list")
        return render_template('index.html', orders=r.json())


if __name__ == '__main__':
    import os

    if 'WINGDB_ACTIVE' in os.environ:
        app.debug = False
    app.run(port=5002)
