# Code Evaluation Exam for Cutler
Welcome
## Running the Web Server
This project runs a web server that handles requests to POST a risk object, get all risk objects, and get a risk object by ID.
To run this service locally, change your working directory to the root of this project and then run `go run main.go`. You can now make requests to the local host on port 8080. If you have another local process running on that port, then you must kill it first.

Once this project is running, you can cURL to the local host endpoints, for example:

```
curl -i -X POST -d '{\"state\": \"open\", \"title\":\"foo\", \"description\":\"blah\"}' http://localhost:8080/v1/risks
curl -i -X POST -d '{\"state\": \"closed\", \"title\":\"bar\", \"description\":\"blah\"}' http://localhost:8080/v1/risks
curl -i -X GET http://localhost:8080/v1/risks
curl -i -X GET http://localhost:8080/v1/risks/1756ca5b-41d9-4e64-9b93-9036579ae15c
```
## Running the Test Suite
There are test suites covering the functional code in this project. These tests leverage the onsi/gomega and onsi/ginkgo pacakges for ease of writing and reading. To run all of the tests, change your working directory to the root of this project and run `go test ./...`. To take full advantage of the onsi/ginkgo output, you can also install ginkgo `go install github.com/onsi/ginkgo/v2/ginkgo` and then run `ginkgo ./...`.