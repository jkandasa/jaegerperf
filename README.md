# Jaeger Performance Tool

### Docker
```
docker run --rm -d -p 8080:8080 --name=jaegerperf quay.io/jkandasa/jaegerperf:1.0
```

### OpenShift
```
oc create -f docker/openshift.yaml
```
