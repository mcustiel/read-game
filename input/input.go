package input

const CODE_EXIT byte = 0
const CODE_JUMP byte = 1
const CODE_RIGHT byte = 2
const CODE_LEFT byte = 3
const CODE_INVALID byte = 4

type Input struct {
	invalid bool
	exit    bool
}

type InputStateAccessor interface {
	IsExit() bool
	IsValid() bool
}

func (inputData Input) IsExit() bool {
	return inputData.exit
}

func (inputData Input) IsValid() bool {
	return !inputData.invalid
}

func NewInput() *Input {
	input := new(Input)
	input.invalid = false
	input.exit = false
	return input
}
