package anet

import (
	"net"
)

var netDriver NetDriver

type NetDriver interface {
	Interfaces() ([]net.Interface, error)
	InterfaceAddrs() ([]net.Addr, error)
}

func SetNetDriver(driver NetDriver) {
	netDriver = driver
}

func Interfaces() ([]net.Interface, error) {
	if netDriver == nil {
		return net.Interfaces()
	}
	return netDriver.Interfaces()
}

func InterfaceAddrs() ([]net.Addr, error) {
	if netDriver == nil {
		return net.InterfaceAddrs()
	}
	return netDriver.InterfaceAddrs()
}
