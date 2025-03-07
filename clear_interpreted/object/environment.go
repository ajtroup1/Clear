package object

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	m := make(map[string]map[string]*Builtin)
	return &Environment{store: s, outer: nil, Modules: m}
}

type Environment struct {
	store   map[string]Object
	outer   *Environment
	Modules map[string]map[string]*Builtin
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

func (e *Environment) GetModule(name string) (map[string]*Builtin, bool) {
	obj, ok := e.Modules[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.GetModule(name)
	}
	return obj, ok
}

func (e *Environment) SetModule(name string, val map[string]*Builtin) {
	e.Modules[name] = val
}
