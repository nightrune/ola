/**
 *
 *
 *
 */
package ola

type MethodDescriptor struct {
	_index       uint32
	_name        string
	_output_type string
	_input_type  string
}

func (m MethodDescriptor) Index() uint32 {
	return m._index
}

func (m MethodDescriptor) Name() string {
	return m._name
}

func (m MethodDescriptor) OutputType() string {
	return m._output_type
}

func (m MethodDescriptor) InputType() string {
	return m._input_type
}

func (m MethodDescriptor) String() string {
	return m._name
}

func NewMethodDescriptor(index uint32, name string, output_type string,
	input_type string) *MethodDescriptor {
	md := new(MethodDescriptor)
	md._index = index
	md._name = name
	md._output_type = output_type
	md._input_type = input_type
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
