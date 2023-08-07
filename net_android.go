//go:build android
// +build android

package anet

import (
	"fmt"
	"net"
	"strings"
)

var netDriver NetDriver

const (
	NetFlagUp           int = iota // interface is up
	NetFlagBroadcast               // interface supports broadcast access capability
	NetFlagLoopback                // interface is a loopback interface
	NetFlagPointToPoint            // interface belongs to a point-to-point link
	NetFlagMulticast               // interface supports multicast access capability
)

type androidNet struct {
	nativeNetDriver NativeNetDriver
}

type NativeNetDriver interface {
	InterfaceAddrs() (*NetInterfaceAddrs, error)
	Interfaces() (*NetInterfaces, error)
}

type NetInterfaceAddrs struct {
	addrs []string
}

type NetInterfaces struct {
	ifaces []*NetInterface
}

type NetInterface struct {
	Index int                // positive integer that starts at one, zero is never used
	MTU   int                // maximum transmission unit
	Name  string             // e.g., "en0", "lo0", "eth0.100"
	Addrs *NetInterfaceAddrs // InterfaceAddresses

	hardwareaddr []byte    // IEEE MAC-48, EUI-48 and EUI-64 form
	flags        net.Flags // e.g., FlagUp, FlagLoopback, FlagMulticast
}

func SetNativeNetDriver(driver NativeNetDriver) {
	netDriver = &androidNet{
		nativeNetDriver: driver,
	}
}

func GetNetDriver() NetDriver {
	return netDriver
}

func (a *androidNet) Interfaces() ([]net.Interface, error) {
	ni, err := a.nativeNetDriver.Interfaces()
	if err != nil {
		return nil, err
	}

	return ni.Interfaces(), nil
}

func (a *androidNet) InterfaceAddrs() ([]net.Addr, error) {
	na, err := a.nativeNetDriver.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	addrs := []net.Addr{}
	for _, addr := range na.addrs {
		if addr == "" {
			continue
		}

		// skip interface name
		ips := strings.Split(addr, "%")
		if len(ips) == 0 {
			continue
		}
		ip := ips[0]

		// resolve ip
		v, err := net.ResolveIPAddr("ip", ip)
		if err != nil {
			continue
		}

		addrs = append(addrs, v)
	}

	return addrs, nil
}

func NewNetInterfaceAddrs() *NetInterfaceAddrs {
	return &NetInterfaceAddrs{addrs: []string{}}
}

func (n *NetInterfaceAddrs) AppendAddr(addr string) {
	n.addrs = append(n.addrs, addr)
}

func (n *NetInterfaces) Interfaces() []net.Interface {
	ifaces := make([]net.Interface, len(n.ifaces))
	for i, iface := range n.ifaces {
		ifaces[i] = iface.Interface()
	}
	return ifaces
}

func (n *NetInterfaces) Append(i *NetInterface) {
	n.ifaces = append(n.ifaces, i)
}

func (n *NetInterface) CopyHardwareAddr(addr []byte) {
	n.hardwareaddr = make([]byte, len(n.hardwareaddr))
	copy(n.hardwareaddr, addr)
}

func (n *NetInterface) Interface() net.Interface {
	return net.Interface{
		Index:        n.Index,
		MTU:          n.MTU,
		Name:         n.Name,
		HardwareAddr: n.hardwareaddr,
		Flags:        n.flags,
	}
}

func (n *NetInterface) AddFlag(flag int) (err error) {
	switch flag {
	case NetFlagUp:
		n.flags |= net.FlagUp
	case NetFlagBroadcast:
		n.flags |= net.FlagBroadcast
	case NetFlagLoopback:
		n.flags |= net.FlagLoopback
	case NetFlagPointToPoint:
		n.flags |= net.FlagPointToPoint
	case NetFlagMulticast:
		n.flags |= net.FlagMulticast
	default:
		err = fmt.Errorf("failed to add unknown flag to net interface: %d", flag)
	}

	return
}
