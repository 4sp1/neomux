package adapter

type Process struct {
	PID    int
	Binary string
}

type Adapter interface {
	List() ([]Process, error)
}
