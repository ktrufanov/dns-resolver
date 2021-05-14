# dns-resolver

An example of a service for working with an k8s-calico-networksets-controller<br>
Resolve DNS names to json(Nets)

Request to dns-resolver:
```
curl -X POST --data "item=github.com"  http://dns-resolver.calico-networksets-controller:8080/dns
```

Response:
```
{"Nets":["140.82.121.4/32"]}
```
