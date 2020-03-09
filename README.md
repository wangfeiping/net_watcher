# net_watcher

## Build

$ git clone https://github.com/wangfeiping/net_watcher.git

$ cd github.com/wangfeiping/net_watcher/

$ sh build.sh

$ cd ./build/

## Try

$ ./net_watcher -h

$ time ./net_watcher call -u https://bing.com

## Deploy

$ ./net_watcher add -u https://bing.com

$ ./net_watcher add -u https://bing.c

$ ./net_watcher start

## Metrics

$ curl http://127.0.0.1:9900/metrics

```
# HELP network_service_status Status of network service response 
# TYPE network_service_status gauge
network_service_status{code="0",url="https://bing.c"} 0
network_service_status{code="200",url="https://bing.com"} 331
```

value = 0 : network service call failed

value > 0 : network service call success, the value is the time(milliseconds) it takes for the service to call

## The configuration of prometheus you may need

prometheus.yml

```
  - job_name: net-watcher
    scrape_interval: 1m
    scrape_timeout: 30s
    static_configs:
      - targets: ['127.0.0.1:9900']
        labels:
          instance: net-watcher
```

rules.yml

```
- name: UrlAccessible
  rules:
  - alert: UrlNotAccessible
    expr: count_over_time(network_service_status{code!="200"}[15m]) > 10
    for: 2m
    labels:
      team: net-watcher
    annotations:
      summary: "Url is not accessible!"
      description: "{{$labels.instance}}: Url ({{$labels.url}}) is not accessible {{ $value }} times in 15m"
  - alert: UrlCallOutOfTime
    expr: network_service_status > 5000
    for: 2m
    labels:
      team: net-watcher
    annotations:
      summary: "Url call out of time"
      description: "Url calling cost {{ $value }} milliseconds: {{$labels.url}}"
```
