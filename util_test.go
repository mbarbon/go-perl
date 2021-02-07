package perl

func eval(i *Interpreter, code string) {
	err := i.EvalVoid(code)
	errPanic(err)
}

func errPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func argTypeMap(i *Interpreter, goValue interface{}) *Scalar {
	arg, err := toPerlArgScalar(i, goValue)
	errPanic(err)
	return newScalarFromMortal(i, arg)
}
