# net_watcher

### build

$ git clone https://github.com/wangfeiping/net_watcher.git

$ cd github.com/wangfeiping/net_watcher/

$ sh build.sh

$ cd ./build/

$ ./net_watcher -h

$ time ./net_watcher call -u https://bing.com

### deploy

$ ./net_watcher add -u https://bing.com

$ ./net_watcher add -u https://bing.com

$ ./net_watcher start

### metrics

$ curl http://127.0.0.1:9900/metrics

```
\# HELP network_service_status Status code of network service response 
\# TYPE network_service_status gauge
network_service_status{code="0",url="https://bing.c"} 0
network_service_status{code="200",url="https://bing.com"} 331

```

value = 0 : network service call failed

value > 0 : cost milliseconds of calling network service
