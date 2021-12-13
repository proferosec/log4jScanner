# log4jScanner

## Goals

This tool provides you with the ability to scan internal (only) subnets for vulnerable log4j web services. It will try to send a JNDI payload to each one (via a User-Agent string and a HTTP header) to a list of common HTTP/S ports. 
For every response it receives, it will log the responding host IP so we can get a list of the vulnerable servers. 

## Basic usage
Download the tool for your specific platform (Windows, Linux or Mac) from the [release](https://github.com/proferosec/log4jScanner/releases/tag/latest) page.
To run the tool, make sure port 5555 on the host is available, and specify the subnet to scan (it is possible to configure a separate address using the `--server` flag):


`
log4jscanner.exe scan --cidr 192.168.7.0/24
`


This will test the top10 HTTP\S ports on the hosts in the subnet and print the vulnerable hosts to the screen and will generate a log in the same location as the binary including all the attempts (both vulnerable and non-vulnerable).
In order to identify which hosts are vulnerable just lookup the word `SUCCESS` in the log.

## Additional usage options
You can use the tool to test for the top 100 HTTP\S ports using the `ports --top100` flag, or for the entire port range using `ports --slow` - Keep in mind, using `ports --slow` will take time to complete.


```
log4jscanner.exe scan --cidr 192.168.7.0/24 --ports top100

log4jscanner.exe scan --cidr 192.168.7.0/24 --ports slow
```


## test setup
Run the docker compose in [here](https://github.com/proferosec/log4jScanner/tree/main/docker):

`docker-compose up -d`

This will provide you with a container vulnerable on port 8080 for HTTP and port 8443 for HTTPS.

Alternativley, you can also run this:
1. Vuln. target: 
   1. `docker run --rm --name vulnerable-app -p 8080:8080 ghcr.io/christophetd/log4shell-vulnerable-app`
2. spin a server for incoming requests
   1. `log4jScanner scanip --cidr DOCKER-RANGE`
3. send a request to the target, with the server details
   1. sends a request to the vuln. target, with the callback details of the sever
   2. once gets a callback, logs the ip of the calling request


### Tests
* test against different subnet
* test for untrusted ssl certs
