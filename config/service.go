package config

import (
	"github.com/micro/go-micro/client/selector"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
)

// 配置micro框架是consul
func RegistryConsul() registry.Registry {
	return consul.NewRegistry(registry.Addrs("172.17.0.5:8500"))
}

// RegistryClient : 注册中心client
func RegistryClient(r registry.Registry) selector.Selector {
	return selector.NewSelector(
		selector.Registry(r),                      //传入consul注册
		selector.SetStrategy(selector.RoundRobin), //指定查询机制
	)
}
