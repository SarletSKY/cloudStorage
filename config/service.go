package config

import (
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
)

// 配置micro框架是consul
func RegistryConsul() registry.Registry {
	return consul.NewRegistry(registry.Addrs("172.17.0.5:8500"))
}
