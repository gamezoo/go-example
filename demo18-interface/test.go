package demo18_interface

type ISystem interface {
	Name() string
}

type SystemBase struct {
}

func (s *SystemBase) Name() string {
	return ""
}

type AwakeSystem struct {
	SystemBase
}
