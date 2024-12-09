package netwatch

import (
	"fmt"
	"log"
	"net"
)

// NetAddr represents a network interface with an index, a name and an IP
type NetAddr struct {
	Index int
	Name  string
	IP    net.IP
}

// Gets the net.IP as IPv4 of a specified network interface
func GetInterfaceIPv4Addr(ifi net.Interface) (net.IP, error) {
	addrs, err := ifi.Addrs()
	if err != nil {
		return nil, err
	}
	var ipv4 net.IP
	for _, addr := range addrs {
		ipv4 = addr.(*net.IPNet).IP.To4()
		if ipv4 != nil {
			return ipv4, nil
		}
	}
	return nil, fmt.Errorf("interface %s has no ipv4 address\n", ifi.Name)
}

func checkFlag(ni net.Interface, f net.Flags) bool {
	if ni.Flags&f == f {
		return true
	}
	return false
}

// Returns the number of network interfaces available on the machine
func InterfaceCount() (int, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return 0, err
	}
	return len(ifaces), nil
}

// Returns a list of all network interfaces as []NetAddr without any filtering
func AllInterfaces() ([]NetAddr, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	na := make([]NetAddr, 0, len(ifaces))
	for _, i := range ifaces {
		ip, err := GetInterfaceIPv4Addr(i)
		if err != nil {
			log.Println("Error getting address for interface:", err)
			ip = net.IPv4(0, 0, 0, 0)
		}
		na = append(na, NetAddr{i.Index, i.Name, ip})
	}
	return na, nil
}

// Returns a list of network interfaces as []NetAddr filtering them by includeFlags and excludeFlags
//
// a common includeFlags is net.FlagRunning to only get the enabled interfaces,
// and a common excludeFlags is net.FlagLoopback to remove the loopback interface from the list
func InterfaceList(includeFlags, excludeFlags net.Flags) ([]NetAddr, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Println("Error getting interface list:", err)
		return nil, err
	}
	na := make([]NetAddr, 0, len(ifaces))
	for _, i := range ifaces {
		if checkFlag(i, excludeFlags) {
			continue
		}
		if !checkFlag(i, includeFlags) {
			continue
		}
		ip, err := GetInterfaceIPv4Addr(i)
		if err != nil {
			log.Println("Error getting address for interface:", err)
			continue
		}
		na = append(na, NetAddr{i.Index, i.Name, ip})
	}
	return na, nil
}
