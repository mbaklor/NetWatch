package net

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// AddressFamily enumeration specifies protocol family and is one of the windows.AF_* constants.
type AddressFamily uint16

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

func notifyIPInterfaceChange(family AddressFamily, callback uintptr, callerContext uintptr, initialNotification bool, notificationHandle *windows.Handle) error {
	var _p0 uint32
	if initialNotification {
		_p0 = 1
	}
	r0, _, _ := syscall.SyscallN(procNotifyIpInterfaceChange.Addr(), uintptr(family), uintptr(callback), uintptr(callerContext), uintptr(_p0), uintptr(unsafe.Pointer(notificationHandle)))
	if r0 != 0 {
		return syscall.Errno(r0)
	}
	return nil
}

func notifyUnicastIpAddressChange(family AddressFamily, callback uintptr, callerContext uintptr, initialNotification bool, notificationHandle *windows.Handle) error {
	var _p0 uint32
	if initialNotification {
		_p0 = 1
	}
	r0, _, _ := syscall.SyscallN(procNotifyUnicastIpAddressChange.Addr(), uintptr(family), uintptr(callback), uintptr(callerContext), uintptr(_p0), uintptr(unsafe.Pointer(notificationHandle)))
	if r0 != 0 {
		return syscall.Errno(r0)
	}
	return nil
}

type NetMonitor struct {
	MonitorChan     chan struct{}
	interfaceHandle *windows.Handle
	addrHandle      *windows.Handle
}

func NewNetMonitor() *NetMonitor {
	ch := make(chan struct{})
	interfaceChange := windows.Handle(0)
	addrChange := windows.Handle(0)
	return &NetMonitor{ch, &interfaceChange, &addrChange}
}

func (nm *NetMonitor) Register() error {
	err := notifyIPInterfaceChange(windows.AF_INET, windows.NewCallback(nm.callback), 0, false, nm.interfaceHandle)
	if err != nil {
		return err
	}
	err = notifyUnicastIpAddressChange(windows.AF_INET, windows.NewCallback(nm.callback), 0, false, nm.addrHandle)
	if err != nil {
		return err
	}
	return nil
}

func (nm *NetMonitor) Unregister() error {
	err := cancelMibChangeNotify2(*nm.interfaceHandle)
	if err != nil {
		return err
	}
	err = cancelMibChangeNotify2(*nm.addrHandle)
	if err != nil {
		return err
	}
	return err
}

func (nm *NetMonitor) callback(callerContext, row, notificationType uintptr) uintptr {
	// fmt.Printf("callback invoked by Windows API (%#v %#v %#v)\n", callerContext, row, notificationType)
	nm.MonitorChan <- struct{}{}
	return 0
}
