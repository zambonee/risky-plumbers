# Comments
* The tests were built off the `onsi/ginkgo` and `onsi/gomega` packages. I have found these to be easier to write, read, and maintain than the standard `testing` package. I hope they help with readability for your purposes too.
* I recommend a graceful shutdown of the service from main.go. I excluded this because it seemed out of scope for this project delivery.
* There are mocks generated by the `mock/mockgen` tool. This made generating mocks easy, with limited code outside of the generated file, to isolate the unit test dependencies.
* I used `gorilla/mux` routing because I find it easier to write and read than the standard `http` package, especially with versioned routing paths.
* There are DAO ("data access object") interface methods that return an error even though the `LocalCache`  of it should not return errors. This is mostly out of habit, especially since we are building an interface that can have an external datasource implementation, which should expect possible errors. Same goes with the `context.Context` parameters.