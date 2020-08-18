import itertools
import time
from multiprocessing import Process, Queue, Value
import uuid
import json

from flask import Flask, g, request
# from tinydb import TinyDB, Query
from kb import kb_engine

app = Flask(__name__)


def kb_run(kill_switch, request_queue, response_queue):
    kb = kb_engine()
    while kill_switch.value == 0:
        request = request_queue.get()
        if request is None:
            time.sleep(0.1)
            continue

        sender_id, message = request
        print(kb, type(kb))
        if sender_id is None:
            kb.tell(message)
        else:
            answer = list(kb.fol_fc_ask(message))
            response_queue.put((sender_id, answer))
        print(kb.clauses)
        print('received:', sender_id, message)
    print('kb engine terminated!')


@app.route('/api/tell', methods=['GET', 'POST'])
def tell():
    if request.method == 'POST':
        query = request.form
    else:
        query = request.args
    request_queue.put((None, query['str']))
    return "success"


@app.route('/api/ask', methods=['GET', 'POST'])
def ask():
    if request.method == 'POST':
        query = request.form
    else:
        query = request.args

    sender_id = str(uuid.uuid4())
    request_queue.put((sender_id, query['str']))

    answer = 'timed out'
    for i in range(5):
        time.sleep(1)
        receiver_id, answer = response_queue.get()
        assert sender_id == receiver_id
        if answer is not None:
            break

    output = {}
    for d in answer:
        for key, value in d.items():
            key = str(key)
            value = str(value)
            if key not in output:
                output[key] = []
            output[key].append(value)
    return json.dumps(output)

request_queue = Queue()
response_queue = Queue()
kill_switch = Value('i', 0)

if __name__ == '__main__':
    engine = Process(target=kb_run, args=(kill_switch, request_queue, response_queue))
    engine.start()
    try:
        app.run()
    except Exception:
        pass
    except KeyboardInterrupt:
        pass
    finally:
        kill_switch.value = 1
        engine.join()
