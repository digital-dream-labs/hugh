package server

// State tracks the state of the connection
type State int

const (
	// Init is a new connection
	Init State = iota
	// Starting means the connection is starting
	Starting
	// Ready means the connection is ready
	Ready
	// Stopping means the connection is being shut down
	Stopping
	// Stopped means the connection is stopping
	Stopped
	// Error means the connection is in error
	Error
)

func (s State) String() string {
	switch s {
	case Init:
		return "INIT"
	case Starting:
		return "STARTING"
	case Ready:
		return "READY"
	case Stopping:
		return "STOPPING"
	case Stopped:
		return "STOPPED"
	case Error:
		return "ERROR"
	default:
		return "INVALID"
	}
}
