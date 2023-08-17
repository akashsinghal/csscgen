# csscgen
Generates kubernetes templates and supply chain artifacts for load testing Container Secure Supply Chain scenarios

## Usage

### Generate K8s resource templates

```bash
# Generate a Deployment template to be used with "clusterloader2" tool
# 1 container, 1 referrer, 1 replica
> csscgen genk8s

apiVersion: apps/v1
kind: Deployment
metadata:
  creationTimestamp: null
  labels:
    group: '{{.Group}}'
  name: '{{.Name}}'
spec:
  replicas: 1
  selector:
    matchLabels:
      name: '{{.Name}}'
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        group: '{{.Group}}'
        name: '{{.Name}}'
    spec:
      containers:
      - image: docker.io/1-containers-1-referrers:1
        name: 1-containers-1-referrers1
        resources: {}
status: {}

```
```bash
# Generate a Deployment template to be used with "clusterloader2" tool
# 2 container, 2 referrer, 10 replica, localhost:5000 registry name
> csscgen genk8s --num-referrers 2 --num-replicas 10 --registry-host localhost:5000 -c 2 

apiVersion: apps/v1
kind: Deployment
metadata:
  creationTimestamp: null
  labels:
    group: '{{.Group}}'
  name: '{{.Name}}'
spec:
  replicas: 10
  selector:
    matchLabels:
      name: '{{.Name}}'
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        group: '{{.Group}}'
        name: '{{.Name}}'
    spec:
      containers:
      - image: localhost:5000/2-containers-2-referrers:1
        name: 2-containers-2-referrers1
        resources: {}
      - image: localhost:5000/2-containers-2-referrers:2
        name: 2-containers-2-referrers2
        resources: {}
status: {}

```

### Generate/Push Images + Referrers
```bash
# The host must be logged into the registry prior to running script
docker login localhost:5000
# Generate 2 unique images + sign each twice and push all artifacts to localhost:5000 registry
./genartifacts.sh -r localhost:5000 -n 2 -s 2
```