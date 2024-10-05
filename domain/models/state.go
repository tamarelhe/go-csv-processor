package models

type State int

const (
	Ready State = iota
	Validated
	Staged
	Processed
	Error
)

func (s State) String() string {
	switch s {
	case Ready:
		return "Ready"
	case Validated:
		return "Validated"
	case Staged:
		return "Staged"
	case Processed:
		return "Processed"
	case Error:
		return "Error"
	default:
		return "Unknown"
	}
}
