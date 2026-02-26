package alert

type Alert struct {
	Type AlertType
	Text string
}

type AlertType int

const (
	Normal AlertType = iota
	Info
	Error
	Warning
)
