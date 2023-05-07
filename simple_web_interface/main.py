# main.py

import requests
from flask import render_template, request, json

from app import app


@app.route('/', methods=['GET', 'POST'])
def index():
    print(request.form)
    if request.form.get("orders") and request.form["send"] == "send":
        r = requests.post(url="http://127.0.0.1:5010/order/get",
                          params={'orderUID': request.form.get("orders")})
        try:
            r.json()
        except:
            pass
            r = requests.get(url="http://127.0.0.1:5010/order/uid/list")
            return render_template('index.html', error = "не найдено заказа с таким UID", orders=r.json())
        pj = str(json.dumps(r.json(), indent=2)).split('\n')

        return render_template('results.html', orders=pj)
    else:
        r = requests.get(url="http://127.0.0.1:5010/order/uid/list")
        return render_template('index.html', orders=r.json())


if __name__ == '__main__':
    import os
    app.run(port=5020)
