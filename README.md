# Jaegerperf Tool

## Description
Jaegerperf tool can be used for the following actions,

* Generate spans
* Measure query execution time

### Generate Spans
Supports two type of spans generation,

* history - generates past data
* realtime - generates future data

to generate spans configuration can be supplied via YAML or JSON format.

#### API: POST http://localhost:8080/api/spansGenerator
Headers: "Content-Type": "application/yaml" or "application/json"
Returns `jobID`.

You can execute via WEB UI as well. 
#### WEB UI: http://localhost:8080 >> Spans Generator

#### Configurations (YAML):
```yaml
target: "collector" # options: agent, collector
endpoint: http://jaegerqe-collector:14268/api/traces
serviceName: jaegerperf_generator
mode: realtime # options: history, realtime
# realtime option (executionDuration)
executionDuration: 5m
# history options (numberOfDays, spansPerDay)
numberOfDays: 10
spansPerDay: 5000
spansPerSecond: 500 # maximum spans limit/sec
childDepth: 4
tags: 
  spans_generator: "jaegerperf"
  days: 10
startTime:
# time format, 2019-01-20T13:34:00
```

* `target` : where do you want to send the spans? to `agent` or `collector`?
* `endpoint` : `agent` or `collector` url
* `serviceName` : service name of the tracer
* `mode` : can be either `realtime` or `history`
* `executionDuration` : this field applicable only if you select the `mode` as `realtime`. can be supplied as `5m`, `1h`, `1h45m` etc.,
* `numberOfDays` : this field applicable only if you select the `mode` as `history`. Number of days do you want to generate the spans
* `spansPerDay` : this field applicable only if you select the `mode` as `history`. Number of spans per day
* `spansPerSecond` : spans limit per second. Sending spans to agent/collector will be limited with this value.
* `childDepth` : Number of child spans should be created
* `tags` : You can define any number of tags. All spans will be containing these tags
* `startTime` : this field is used for `history` mode. if you leave this field current time will be start time.


### Measure query execution time
Measures query execution with given input.

to execute query runner configuration can be supplied via YAML or JSON format.

#### API: POST http://localhost:8080/api/queryRunner
Headers: "Content-Type": "application/yaml" or "application/json"
Returns `jobID`.

You can execute via WEB UI as well. 
#### WEB UI: http://localhost:8080 >> Query Runner

#### Configurations (YAML):
```yaml
hostUrl: http://jaegerqe-query:16686
tags:
  - master branch
tests:
  - name: last 12 hours limit 100
    type: search
    iteration: 10
    statusCode: 200
    queryParams:
      service: jaegerperf_generator
      limit: 100
      lookback: custom
      start: -12h
      end: 0h

  - name: last 7 days limit 2000
    type: search
    iteration: 10
    statusCode: 200
    queryParams:
      service: jaegerperf_generator
      limit: 2000
      lookback: custom
      start: -168h
      end: 0h
```
* `hostUrl` : jaegertracing query service url
* `tags` : supply any number tags to refer this result on the `Query Metrics` page. Unique tags are recommended for each run
* `tests` : you can add list of tests to run
  * `name` : Name of the test
  * `type` : Supports jaeger query `search` and `services` query.
  * `iteration` : Number of time you want to sample the query. Number of time you want to execute the same query
  * `statusCode` : expected status code of the query response
  * `queryParams` : input parameters for the query. You can add any number of parameters, that supports your query. Here is an example for jaeger-query search api
    * `service` : service name
    * `limit` : limit the spans count on the response
    * `start` : supports dynamically. 
      * examples: 
        * `1h` -> current time + 1 hour
        * `-1h` -> current time - 1 hour
        * `0h` -> current time
        * `-12h` -> current time - 12 hours
        * `-168h` -> current time - 168 hours (ie: current time - 7 days)

### Jobs
Once triggered `Query Runner` or `Spans Generator` you can see the status of the the job. You can get exact `jobID` when you trigger the API.

#### API: POST http://localhost:8080/api/jobs

#### WEB UI, http://localhost:8080  >> Jobs
In the WEB UI, for Query runner, you can see a quick summary table on a completed jobs.

#### JOB DATA
You can download job details via the API call or from the UI page. It contains detailed information.
Example: You can see Query Runner data like elapsed time, failed count, etc.,

### Query Metrics (WEB UI)
Based on the tags supplied on the `Query Runner` job you can see metric data in the form of table and charts.

It give quick understanding. If you want to see failed query detail like status code and error message can download from the `Jobs` page.


## RUN
### Docker
```
docker run --rm -d -p 8080:8080 --name=jaegerperf quay.io/jkandasa/jaegerperf:1.1
```

### OpenShift
```
oc create -f docker/openshift.yaml
```

### Kubernetes
```
kubectl create -f docker/kubernetes.yaml
```

## Build
Build WEB
```
cd web
npm run build
```
Build and push docker image
```
docker/build_docker.sh
```