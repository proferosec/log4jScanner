# log4jScanner

## Goals

This tool provides you with the ability to scan internal (only) subnets for vulnerable log4j web services. 
It will attempt to send a JNDI payload to each discovered web service (via the [methods](#methods_used) outlined below) to a list of common HTTP/S ports. 
For every response it receives, it will log the responding host IP so we can get a list of the vulnerable servers.

If there is a "SUCCESS", this means that some web service has received the request, was vulnerable to the **log4j** exploit and sent a request to our TCP server.

The tool does not send any exploits to the vulnerable hosts, and is designed to be as passive as possible.

## Latest Release

| Platform | Binary   | Checksum |
|----------|----------|----------|
| Windows  |[log4jscanner-windows.zip](https://github.com/proferosec/log4jScanner/releases/download/latest/log4jscanner-windows.zip) | [SHA256](https://github.com/proferosec/log4jScanner/releases/download/latest/windows.sha256.txt) |
| Linux  |[log4jscanner-linux.zip](https://github.com/proferosec/log4jScanner/releases/download/latest/log4jscanner-linux.zip) | [SHA256](https://github.com/proferosec/log4jScanner/releases/download/latest/linux.sha256.txt) |
| MacOS  |[log4jscanner-darwin.zip](https://github.com/proferosec/log4jScanner/releases/download/latest/log4jscanner-darwin.zip) | [SHA256](https://github.com/proferosec/log4jScanner/releases/download/latest/darwin.sha256.txt) |

## Example

![example](https://github.com/proferosec/log4jScanner/blob/staging/movie.gif)

In this example we run the tool against the `192.168.1.0/24` subnet


## Basic usage
Download the tool for your specific platform (Windows, Linux or Mac), to run the tool, make sure port 5555 on the host is available (or change it via configuration), 
and specify the subnet to scan (it is possible to configure a separate server:port combination using the `--server` flag):

```bash
log4jScanner.exe scan --cidr 192.168.7.0/24
```

This will test the top 10 HTTP\S ports on the hosts in the subnet,  print any vulnerable hosts to the screen, 
and generate a log + summary CSV in the same location as the binary including all the attempts (both vulnerable and non-vulnerable).

In order to identify which hosts are vulnerable just lookup the word `SUCCESS` in the log, you can grep the log for the keywork `SUCCESS` to get just the results.
Also, the tool generates a CSV file containing all the results, filter on `vulenrable` to get the vulnerable hosts.

### Additional usage options
You can use the tool to test for the top 100 HTTP\S ports using the `ports אop100` flag, or for the entire port range using `ports slow` - Keep in mind, using `ports slow` will take time to complete.

```bash
log4jscanner.exe scan --cidr 192.168.7.0/24 --ports=top100
```

it is possible to use a non-default configuration for the callback server
```bash
log4jscanner.exe scan --cidr 192.168.7.0/24 --server=192.168.1.100:5000
```

if you wish to disable the callback server, use `--noserver`

### Available flags

* `--nocolor` provide output without color
* `--ports` either top10 (default) or top100 (list of the 100 most commong web ports)
* `--noserver` only scan, do not use a local callback server

### Methods Used

Currently the tool uses the following areas to try and send an exploit

### Test setup

In order to test your environment, you can use the included docker images to launch vulnerable applications.

Run the docker compose in [here](https://github.com/proferosec/log4jScanner/tree/main/docker):

`docker-compose up -d`

This will provide you with a container vulnerable on port 8080 for HTTP and port 8443 for HTTPS.

Alternativley, you can also run this:
1. Vuln. target: 
   1. `docker run --rm --name vulnerable-app -p 8080:8080 ghcr.io/christophetd/log4shell-vulnerable-app`
2. spin a server for incoming requests
   1. `log4jScanner scanip --cidr DOCKER-SUBNET`
3. send a request to the target, with the server details
   1. sends a request to the vuln. target, with the callback details of the sever
   2. once gets a callback, logs the ip of the calling request


# Contributions

We welcome contributions, please submit a PR or contact us via contact@profero.io
