# http-load-tester
HTTP load testing in Go.

## Installation
Clone the repo & build the tool:

`git clone https://github.com/Aidan-Simard/http-load-tester.git && cd http-load-tester`

`go build load.go`

## Usage
`./load [options...]`

### Options
-u (required): The URL to test

-n: The number of requests to make (default=1)

-c: The number of concurrent threads to send requests with (default=1)

## Example
`./load -u http://localhost:8081 -n 2000 -c 2`

Send 2000 requests (concurrency of 2) to http://localhost:8081

Output:

```
Successes: 2000
Failures: 0
Total Time: 1296 (ms)
Requests Per Second: 1543

Min Request Time: 300.000000 (us)
Mean Request Time: 1292.935000 (us)
Max Request Time: 84771.000000 (us)
```
