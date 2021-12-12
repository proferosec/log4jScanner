# log4j-scanner

## Goals

this tool will scan a subnet for web servers and will try to send the JDNI exploit to each one. 
For every response it receives, it will log the sender IP so we can get a list of the vulnerable servers. 

## Design

The utility spins up a webserver listening for incoming requests. 
then, it will open a request for every available port in the range that responds to HTTP/S and send it the exploit.

1. get all IPs in the CIDR
2. scan each IP for open ports (either complete list, or reduced list)
3. for any open port, call the `ScanIP` 
4. log all callbacks (source IP address)
5. the callback server is listening to `localhost:5555`
6. if the `--slow` flag is used, all ports are scanned, for each IP

## test setup

1. Vuln. target: 
   1. `docker run --rm --name vulnerable-app -p 8080:8080 ghcr.io/christophetd/log4shell-vulnerable-app`
2. spin a server for incoming requests
   1. `log4jScanner scanip -s --cidr 192.168.1.0/24`
3. send a request to the target, with the server details
   1. sends a request to the vuln. target, with the callback details of the sever
   2. once gets a callback, logs the ip of the calling request


### Tests
* test against different subnet
* test for untrusted ssl certs
