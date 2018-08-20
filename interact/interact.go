package interact

// RunFace for interact methods
type RunFace interface {
	Run() *Value
}

// Interactive
type Interactive struct {
	Name string
}

func New(name string) *Interactive {
	return &Interactive{Name: name}
}

// Option
type Option struct {
	Quit bool
	// default value
	DefVal string
}
