import itertools
import time
from multiprocessing import Process, Queue, Value
import uuid
import json
import logging

from flask import Flask, g, request
# from tinydb import TinyDB, Query
from kb import kb_engine

logging.basicConfig(
    format='%(asctime)s %(levelname)-8s %(message)s',
    level=logging.INFO,
    datefmt='%Y-%m-%d %H:%M:%S')

app = Flask(__name__)


def kb_run(kill_switch, request_queue, response_queue):
    kb = kb_engine()
    logging.info('KB engine started.')
    while kill_switch.value == 0:
        request = request_queue.get()
        if request is None:
            time.sleep(0.1)
            continue

        try:
            sender_id, message = request
            if sender_id is None:
                kb.tell(message)
            else:
                if message == "dump":
                    response_queue.put((sender_id, kb.clauses))
                else:
                    answer = list(kb.fol_fc_ask(message))
                    response_queue.put((sender_id, answer))
        except Exception as ex:
            logging.error('KB engine: {}'.format(str(ex)))
    logging.info('KB engine terminated.')


@app.route('/api/tell', methods=['GET', 'POST'])
def tell():
    if request.method == 'POST':
        query = request.form
    else:
        query = request.args
    request_queue.put((None, query['str']))
    return "success"


@app.route('/api/sense', methods=['GET', 'POST'])
def sense():
    if request.method == 'POST':
        query = request.form
    else:
        query = request.args
    

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
        if answer is not None:
            break

    output = {}
    for d in answer:
        for key, value in d.items():
            key = str(key)
            value = str(value)
            if key not in output:
                output[key] = []
            if value not in output[key]:
                output[key].append(value)
    return json.dumps(output)


@app.route('/api/list_kb', methods=['GET', 'POST'])
def list_kb():
    sender_id = str(uuid.uuid4())
    request_queue.put((sender_id, "dump"))

    answer = 'timed out'
    for i in range(5):
        time.sleep(1)
        receiver_id, answer = response_queue.get()
        if answer is not None:
            break
    output = [str(c) for c in answer]
    return '\n'.join(output)

request_queue = Queue()
response_queue = Queue()
kill_switch = Value('i', 0)

if __name__ == '__main__':
    engine = Process(target=kb_run, args=(kill_switch, request_queue, response_queue))
    engine.start()
    try:
        app.run(host='0.0.0.0')
    except Exception:
        pass
    except KeyboardInterrupt:
        pass
    finally:
        kill_switch.value = 1
        engine.join()
