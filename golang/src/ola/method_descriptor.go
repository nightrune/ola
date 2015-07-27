/**
 *
 */
package ola

type MethodDescriptor struct {
	_index uint32
	_name  string
}

func (m *MethodDescriptor) Index() uint32 {
	return m._index
}

func (m *MethodDescriptor) Name() string {
	return m._name
}

func NewMethodDescriptor(index uint32, name string) *MethodDescriptor {
	md := new(MethodDescriptor)
	md._index = index
	md._name = name
	return md
}

type InvalidMethod struct {
	_s string
}

func (m *InvalidMethod) Error() string {
	return m._s
}

func NewInvalidMethod(err_string string) *InvalidMethod {
	e := new(InvalidMethod)
	e._s = err_string
	return e
}
