package fury

type Fury struct {
	Name string
}

func New() (f *Fury) {
	f = &Fury{
		Name: "Hello Fury!!!",
	}
	return
}
