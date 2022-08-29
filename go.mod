module Balancer

go 1.17

require (
	github.com/gorilla/mux v1.8.0
	github.com/slairmy/balancer/balancer v1.0.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/lafikl/consistent v0.0.0-20220512074542-bdd3606bfc3e // indirect
	github.com/minio/blake2b-simd v0.0.0-20160723061019-3f5f724cb5b1 // indirect
	github.com/nxadm/tail v1.4.8 // indirect
	github.com/starwander/GoFibonacciHeap v0.0.0-20190508061137-ba2e4f01000a // indirect
	golang.org/x/net v0.0.0-20220425223048-2871e0cb64e4 // indirect
	golang.org/x/sys v0.0.0-20220422013727-9388b58f7150 // indirect
)

replace github.com/slairmy/balancer/balancer v1.0.0 => ./balancer
