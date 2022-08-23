module Balancer

go 1.17

require (
	github.com/slairmy/balancer/balancer v1.0.0
	gopkg.in/yaml.v3 v3.0.1
)

require github.com/gorilla/mux v1.8.0 // indirect

replace github.com/slairmy/balancer/balancer v1.0.0 => ./balancer
