# _This repo is archived. See http://github.com/waggle-sensor/edge-scheduler_

## Knowledgebase for Node Scheduler

This knowledgebase (KB) offers an inference engine for logical entailment. The KB accepts local perceptions from sensor readings as well as user-defined rules. Those perceptions must be expressed in [first-order logic](https://en.wikipedia.org/wiki/First-order_logic). Examples are,

```
# it rains in the local region
Raining(local)
# it rains globally, e.g., entire city, state, etc
Raining(global)
# Smoke detector requires GPU
Require(Smoke, GPU)
# it is daytime now
Daytime(Now)
# if more than 2 plugins require GPU, they conflict
Require(x, GPU) & Require(y, GPU) ==> Conflict(x, y)
```

To run KB,
```
# assume the name of Docker image is waggle/knowledgebase
$ docker run -d -p 5000:5000 waggle/knowledgebase
```

To add perceptions,
```
# assume curl is installed on the host
$ curl http://localhost:5000/api/tell?str=Require\(Smoke,GPU\)
success
$ curl http://localhost:5000/api/tell?str=Require\(Cloud,GPU\)
success
$ curl --data-urlencode "str=Require(x, GPU) & Require(y, GPU) ==> Conflict(x, y)" http://localhost:5000/api/tell
success
```

To ask a query,
```
# query if Smoke conflicts with any
$ curl -d "str=Conflict(Smoke, y)" http://localhost:5000/api/ask
{"y": ["Smoke", "Cloud"]}
```

### Build

Use [Dockerfile](Dockerfile) to build a Docker image.

### Acknowledgement

The knowledgebase engine is implemented on the basis of the code provided from the book "Artificial Intelligence: A Modern Approach" written by Stuart J. Russell and Peter Norvig. The code book is hosted in [the Github repository](https://github.com/aimacode)
