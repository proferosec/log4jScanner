# log4jScanner

## Goals

This tool provides you the ability to scan internal (only) subnets for vulnerable log4j services running. It will try to send a JDNI payload to each one (via the User-agent string and a HTTP header. 
For every response it receives, it will log the responding host IP so we can get a list of the vulnerable servers. 

## Basic usage
Download the tool for your specific platform (Windows, Linux or Mac) from the [release](https://github.com/proferosec/log4jScanner/releases/tag/latest) page.
To run the tool, make sure port 5555 on the host is available, and specify the subnet to scan:


`
log4jscanner.exe scan --cidr 192.168.7.0/24
`


This will test the top10 HTTP\S ports on the hosts in the subnet and print the volnurable hosts to the screen and will generate a log in the same location as the binary including all the attempts (both vulnerable and non-volnerable).
In order to identify which hosts are volnerable just lookup the word `SUCCESS` in the log.

## Additional usage options
You can use the tool to test for the top 100 HTTP\S ports using the `--top100` flag, or for the entire port range using `--slow` - Keep in mind, using `--slow` will take time to complete.


## test setup

1. Vuln. target: 
   1. `docker run --rm --name vulnerable-app -p 8080:8080 ghcr.io/christophetd/log4shell-vulnerable-app`
2. spin a server for incoming requests
   1. `log4jScanner scanip -s --cidr DOCKER-RANGE`
3. send a request to the target, with the server details
   1. sends a request to the vuln. target, with the callback details of the sever
   2. once gets a callback, logs the ip of the calling request


### Tests
* test against different subnet
* test for untrusted ssl certs
