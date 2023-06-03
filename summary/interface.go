package summary

// Interface is an interface for summarizing configuration.
type Interface interface {
	Add(format string, a ...interface{})
	Descend()
	Ascend()
}
