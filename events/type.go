package events

type Fetcher interface {
	Fetch(limit int) ([]Event, error)
}

type Processor interface {
	Process(e Event) error
}

type Type int

// iota автоматически инкрементирует значение для груп констант
const (
	Unknown Type = iota
	Message
)

type Event struct {
	Type Type
	Text string
}
