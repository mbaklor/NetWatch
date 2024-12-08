package net

import (
	"fmt"
	"log"
	"net"
)

type NetAddr struct {
	Index int
	Name  string
	IP    string
}

func GetInterfaceIPv4Addr(ifi net.Interface) (string, error) {
	addrs, err := ifi.Addrs()
	if err != nil {
		return "", err
	}
	var ipv4 net.IP
	for _, addr := range addrs {
		ipv4 = addr.(*net.IPNet).IP.To4()
		if ipv4 != nil {
			return ipv4.String(), nil
		}
	}
	return "", fmt.Errorf("interface %s has no ipv4 address\n", ifi.Name)
}

func checkFlag(ni net.Interface, f net.Flags) bool {
	if ni.Flags&f == f {
		return true
	}
	return false
}

func InterfaceCount() (int, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return 0, err
	}
	return len(ifaces), nil
}

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
			ip = "0.0.0.0"
		}
		na = append(na, NetAddr{i.Index, i.Name, ip})
	}
	return na, nil
}

func InterfaceList() ([]NetAddr, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Println("Error getting interface list:", err)
		return nil, err
	}
	na := make([]NetAddr, 0, len(ifaces))
	for _, i := range ifaces {
		if checkFlag(i, net.FlagLoopback) {
			continue
		}
		if !checkFlag(i, net.FlagRunning) {
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
