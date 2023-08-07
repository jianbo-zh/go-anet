package anet

import (
	"net"
)

type NetDriver interface {
	Interfaces() ([]net.Interface, error)
	InterfaceAddrs() ([]net.Addr, error)
}

func Interfaces() ([]net.Interface, error) {
	if GetNetDriver() == nil {
		return net.Interfaces()
	}
	return GetNetDriver().Interfaces()
}

func InterfaceAddrs() ([]net.Addr, error) {
	if GetNetDriver() == nil {
		return net.InterfaceAddrs()
	}
	return GetNetDriver().InterfaceAddrs()
}
