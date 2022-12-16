package monitor

type Monitor interface {
	Start() (<-chan Event, error)
	Stop() error
}

type RecordingStatus byte

const (
	NotRecording RecordingStatus = iota
	Recording
)

type Event struct {
	RecordingStatus RecordingStatus
}
