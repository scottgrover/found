## Found 
### A simple service mDns based service discovery tool written in golang. 
The goal was to create a small binary that I can set up as a service on boxes deployed in my home network, so I can keep track of them if their IP's change. 

### Usage:
The service requries a configuration file. An example configuration can be found in config.yaml. 
```./main -config config.yaml```

### How to build:
```go build main.go```

### How does it work?
There are two modes, discovery and service:
* Service mode: runs on the box that you want to keep track of. 

   * This uses mDns to make the service discoverable. 
   * It has the ability to share it's:
      node name, the type service it hosts, the protocol that service uses, the port for the service, and any metadata you want to configure.
   
* Discovery mode: will find the services on the network placed in service mode.

   * Requires you to know the type of service it's looking for a protocol, and it's domain. 

### Resources:
- https://github.com/grandcat/zeroconf
- https://sosedoff.com/2017/09/07/zeroconf.html
