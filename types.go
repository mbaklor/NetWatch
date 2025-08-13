package netwatch

import "net"

// Indicates what callback was called in the notification
type CallerContext uint8

const (
	CONTEXT_NOTIFY_IP_INTERFACE_CHANGE CallerContext = iota
	CONTEXT_NOTIFY_UNICAST_IP_ADDRESS_CHANGE
)

func (cc CallerContext) String() string {
	switch cc {
	case CONTEXT_NOTIFY_IP_INTERFACE_CHANGE:
		return "IP Interface change"
	case CONTEXT_NOTIFY_UNICAST_IP_ADDRESS_CHANGE:
		return "Unicast IP Address change"
	default:
		return ""
	}
}

// The MIB_NOTIFICATION_TYPE enumeration type defines the notification type that is passed to a callback function when a notification occurs.
type MibNotificationType uint8

const (
	// A parameter was changed.
	MIB_PARAMETER_NOTIFICATION MibNotificationType = iota
	// A new MIB instance was added.
	MIB_ADD_INSTANCE
	// A new MIB instance was added.
	MIB_DELETE_INSTANCE
	// A notification that is invoked immediately after registration for change notification completes.
	// This initial notification does not indicate that a change occurred to a MIB instance.
	// The purpose of this initial notification type is to provide confirmation that the callback function is properly registered.
	MIB_INITIAL_NOTIFICATION
)

func (mnt MibNotificationType) String() string {
	switch mnt {
	case MIB_PARAMETER_NOTIFICATION:
		return "Parameter Notification"
	case MIB_ADD_INSTANCE:
		return "Add Instance"
	case MIB_DELETE_INSTANCE:
		return "Delete Instance"
	case MIB_INITIAL_NOTIFICATION:
		return "Initial Notification"
	default:
		return ""
	}
}

type MonitorNotification struct {
	CallerContext    CallerContext
	NotificationType MibNotificationType
	InterfaceIndex   int
	InterfaceAddress net.IP
}
