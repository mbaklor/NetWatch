package netwatch

import (
	"errors"
	"net"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	modiphlpapi = windows.NewLazySystemDLL("iphlpapi.dll")

	procCancelMibChangeNotify2       = modiphlpapi.NewProc("CancelMibChangeNotify2")
	procNotifyUnicastIpAddressChange = modiphlpapi.NewProc("NotifyUnicastIpAddressChange")
	procNotifyIpInterfaceChange      = modiphlpapi.NewProc("NotifyIpInterfaceChange")
)

func cancelMibChangeNotify2(handler windows.Handle) error {
	r0, _, _ := syscall.SyscallN(procCancelMibChangeNotify2.Addr(), uintptr(handler))
	if r0 != 0 {
		return syscall.Errno(r0)
	}
	return nil
}

func notifyIPInterfaceChange(family uintptr, callback uintptr, callerContext CallerContext, initialNotification bool, notificationHandle *windows.Handle) error {
	var _p0 uint32
	if initialNotification {
		_p0 = 1
	}
	r0, _, _ := syscall.SyscallN(procNotifyIpInterfaceChange.Addr(), family, callback, uintptr(callerContext), uintptr(_p0), uintptr(unsafe.Pointer(notificationHandle)))
	if r0 != 0 {
		return syscall.Errno(r0)
	}
	return nil
}

func notifyUnicastIpAddressChange(family uintptr, callback uintptr, callerContext CallerContext, initialNotification bool, notificationHandle *windows.Handle) error {
	var _p0 uint32
	if initialNotification {
		_p0 = 1
	}
	r0, _, _ := syscall.SyscallN(procNotifyUnicastIpAddressChange.Addr(), family, callback, uintptr(callerContext), uintptr(_p0), uintptr(unsafe.Pointer(notificationHandle)))
	if r0 != 0 {
		return syscall.Errno(r0)
	}
	return nil
}

// NetMonitor represents a registered monitor to the network interfaces on windows machines
// once registered, any change to a network name or ip address is sent to the MonitorChan
type NetMonitor struct {
	MonitorNotificationChan chan MonitorNotification
	interfaceHandle         *windows.Handle
	addrHandle              *windows.Handle
}

// Creates a new NetMonitor and initiates its channel
func NewNetMonitor() *NetMonitor {
	ch := make(chan MonitorNotification)
	interfaceChange := windows.Handle(0)
	addrChange := windows.Handle(0)
	return &NetMonitor{ch, &interfaceChange, &addrChange}
}

// Registers the network monitor to the windows notifier
func (nm *NetMonitor) Register(initialNotification bool) error {
	err1 := notifyIPInterfaceChange(windows.AF_INET, windows.NewCallback(nm.callbackIpInterfaceChange), CONTEXT_NOTIFY_IP_INTERFACE_CHANGE, initialNotification, nm.interfaceHandle)
	err2 := notifyUnicastIpAddressChange(windows.AF_INET, windows.NewCallback(nm.callbackUnicastAddressChange), CONTEXT_NOTIFY_UNICAST_IP_ADDRESS_CHANGE, initialNotification, nm.addrHandle)
	return errors.Join(err1, err2)
}

// Removes the windows notifier to the network interfaces
func (nm *NetMonitor) Unregister() error {
	var err1, err2 error
	if nm.interfaceHandle != nil {
		err1 = cancelMibChangeNotify2(*nm.interfaceHandle)
	}
	if nm.addrHandle != nil {
		err2 = cancelMibChangeNotify2(*nm.addrHandle)
	}
	return errors.Join(err1, err2)
}

func sockaddrToIP(rsa windows.RawSockaddrInet6) net.IP {
	switch rsa.Family {
	case windows.AF_INET:
		rsa4 := (*windows.RawSockaddrInet4)(unsafe.Pointer(&rsa))
		return net.IP(rsa4.Addr[:])
	case windows.AF_INET6:
		return net.IP(rsa.Addr[:])
	default:
		return nil
	}
}

func (nm *NetMonitor) callbackIpInterfaceChange(callerContext CallerContext, row *windows.MibIpInterfaceRow, notificationType MibNotificationType) uintptr {
	nm.MonitorNotificationChan <- MonitorNotification{
		InterfaceIndex:   int(row.InterfaceIndex),
		NotificationType: (notificationType),
		CallerContext:    (callerContext),
		InterfaceAddress: nil,
	}
	return 0
}
func (nm *NetMonitor) callbackUnicastAddressChange(callerContext CallerContext, row *windows.MibUnicastIpAddressRow, notificationType MibNotificationType) uintptr {
	nm.MonitorNotificationChan <- MonitorNotification{
		InterfaceIndex:   int(row.InterfaceIndex),
		NotificationType: (notificationType),
		CallerContext:    (callerContext),
		InterfaceAddress: sockaddrToIP(row.Address),
	}
	return 0
}
