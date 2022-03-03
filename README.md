# vmware-work-sample
Sample GoLang service that reports prometheus metrics

[![test](https://github.com/jtcressy/vmware-work-sample/actions/workflows/test.yml/badge.svg)](https://github.com/jtcressy/vmware-work-sample/actions/workflows/test.yml)
[![release](https://github.com/jtcressy/vmware-work-sample/actions/workflows/release.yml/badge.svg)](https://github.com/jtcressy/vmware-work-sample/actions/workflows/release.yml)


# Overview
vmware-work sample is a small service to demonstrate my skills and ability to create a golang application that expresses metrics and can live within a kubernetes cluster.

The basic job of this service is to query some URLs and log their response time and status. By default it queries `https://httpstat.us/200` and `https://httpstat.us/503`.

The service exposes two endpoints:
- /metrics
- /status

Most important of which is `/metrics` which returns prometheus-formatted metric data. Aside from library default metrics, the two dimensions added by this service are:
  - `sample_external_url_up`
    - Boolean status of whether a URL is considered up or down.
    - Labels: `url`
  - `sample_external_url_response_ms`
    - Response time in milliseconds it took for the URL to respond.
    - Labels: `url`

example:
```
❯ curl -s "localhost:8080/metrics" | tail
promhttp_metric_handler_requests_total{code="500"} 0
promhttp_metric_handler_requests_total{code="503"} 0
# HELP sample_external_url_response_ms Response time in milliseconds it took for the URL to respond.
# TYPE sample_external_url_response_ms gauge
sample_external_url_response_ms{url="https://httpstat.us/200"} 96
sample_external_url_response_ms{url="https://httpstat.us/503"} 145
# HELP sample_external_url_up Boolean status of whether a URL is considered up or down.
# TYPE sample_external_url_up gauge
sample_external_url_up{url="https://httpstat.us/200"} 1
sample_external_url_up{url="https://httpstat.us/503"} 0

```

The `/status` endpoint simply returns a json document with information about the queried URLs such as up status and response time in milliseconds.

example:
```
❯ curl "localhost:8080/status?pretty=true"
{
    "https://httpstat.us/200": {
        "response_time_ms": 107,
        "up": true
    },
    "https://httpstat.us/503": {
        "response_time_ms": 131,
        "up": false
    }
}
```

## Sample Dashboard

A sample grafana dashboard can be found at `/grafana/dashboards`

Additionally, screenshots of prometheus and grafana can be seen at `/screenshots`

# Usage

## Container image

A container image for this project is continuously built by github actions and is available at the following tag:
```
ghcr.io/jtcressy/vmware-work-sample:latest
```
Example docker command & usage:
```
docker run --rm -p 8080:8080 -d ghcr.io/jtcressy/vmware-work-sample:latest -bind-addr 0.0.0.0:8080 -test-url https://google.com

curl "localhost:8080/status?pretty"
```

## CLI Arguments

```
❯ vmware-work-sample -h
Usage of vmware-work-sample:
  -bind-addr string
        The address the webserver binds to. (default ":8080")
  -fetch-interval duration
        How often to ping URL endpoints (default: 30s) (default 5s)
  -test-url value
        A URL to query for uptime statistics. Use multiple times to query multple URL's in parallel
```

Example using all arguments:
```
❯ ./vmware-work-sample -bind-addr 0.0.0.0:8080 -fetch-interval 10s -test-url https://httpstat.us/201 -test-url https://httpstat.us/404
2022/03/02 20:07:42 Starting pinger with status server on 0.0.0.0:8080
2022/03/02 20:07:52 Running ping for https://httpstat.us/201
2022/03/02 20:07:52 Running ping for https://httpstat.us/404
```

## Building binary from source

To build vmware-work-sample from the source code yourself, you need to have a working Go environment with version 1.17 or greater installed.

```
git clone https://github.com/jtcressy/vmware-work-sample.git
cd vmware-work-sample
make build
```

## Deployment

Kustomize configurations are available under `/config` with the following overlays:
- `/config/overlays/default`
  - Deploys the most recent tagged image from ghcr.io
- `/config/overlays/google-pinger`
  - Deploys the most recent tagged image from ghcr.io
  - Configures deployment to ping https://google.com instead of the defaults.
- `/config/overlays/local-minikube`
  - Deploys image tagged with ghcr.io/jtcressy/vmware-work-sample:dev
  - Intended to be used with local minikube environment where docker image builds happen

e.g.
```
kubectl apply -k https://github.com/jtcressy/vmware-work-sample/config/overlays/default?ref=main
```

# Local Development

## Testing

Simple make directive:
```
make test
```

## Deploying Locally with Minikube

***Ensure you have docker, minikube and kubectl running on your environment.***

```
git clone https://github.com/jtcressy/vmware-work-sample.git
cd vmware-work-sample
make deploy-local
```

Simply rerun `make deploy-local` to pick up new changes and deploy them.