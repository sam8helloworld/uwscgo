package object

type BindedObjectType string

const (
	LOCAL  = "LOCAL"
	PUBLIC = "PUBLIC"
	CONST  = "CONST" // 上書き不可
)

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

type BindedObject struct {
	Key    string // なくてもいいけどデバッグ用に入れている
	Object Object
	Type   BindedObjectType
}

type Environment struct {
	store map[string]*BindedObject
	outer *Environment
}

func NewEnvironment() *Environment {
	s := make(map[string]*BindedObject)
	return &Environment{
		store: s,
		outer: nil,
	}
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.BindedObject(name)
	}
	if obj != nil {
		return obj.Object, ok
	} else {
		return nil, ok
	}
}

func (e *Environment) BindedObject(name string) (*BindedObject, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.BindedObject(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = &BindedObject{
		Key:    name,
		Object: val,
		Type:   LOCAL,
	}
	return val
}

func (e *Environment) SetConst(name string, val Object) Object {
	e.store[name] = &BindedObject{
		Key:    name,
		Object: val,
		Type:   CONST,
	}
	return val
}
