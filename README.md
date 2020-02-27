# net_watcher

\# HELP network_service_status Status code of network service response 
\# TYPE network_service_status gauge
network_service_status{code="0",url="https://bing.c"} 0
network_service_status{code="200",url="https://bing.com"} 331

value = 0 : network service call failed
value > 0 : cost milliseconds of calling network service
