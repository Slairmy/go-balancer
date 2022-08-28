module Balancer

go 1.17

require (
	github.com/gorilla/mux v1.8.0
	github.com/slairmy/balancer/balancer v1.0.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.20.1 // indirect
	github.com/starwander/GoFibonacciHeap v0.0.0-20190508061137-ba2e4f01000a // indirect
)

replace github.com/slairmy/balancer/balancer v1.0.0 => ./balancer
