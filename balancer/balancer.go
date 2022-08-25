package balancer

import "errors"

var (
	AlgorithmNotSupportedError = errors.New("algorithm not supported")
)

// 负载均衡接口定义 -- 添加Add，移除Remove，Balance选择主机，主机连接数+1 Inc，主机连接数 -1
// Add 和 Remove 主要是为了在健康检查的时候将无效的host或者有效的host加进来

type Balancer interface {
	Add(string)
	Remove(string)
	Balance(string) (string, error)
	Inc(string)
	Done(string)
}

// Factory 工厂模式 -- 算法生成器
type Factory func([]string) Balancer

var factories = make(map[string]Factory)

// Build 负载均衡策略 -- 传入负载均衡算法类型参数
func Build(algorithm string, hosts []string) (Balancer, error) {
	factory, ok := factories[algorithm]

	if !ok {
		return nil, AlgorithmNotSupportedError
	}

	return factory(hosts), nil
}
