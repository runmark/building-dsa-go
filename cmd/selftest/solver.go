package selftest

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
)

type Processor struct {
	Solver MathSolver
}

func (p Processor) ProcessExpression(ctx context.Context, r io.Reader) (float64, error) {
	curExpression, err := readToNewLine(r)
	if err != nil {
		return 0, nil
	}

	if len(curExpression) == 0 {
		return 0, errors.New("no expression to read")
	}

	answer, err := p.Solver.Resolve(ctx, curExpression)
	return answer, err
}

func readToNewLine(r io.Reader) (string, error) {
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return string(bytes), err
	}
	return string(bytes), err
}

type MathSolver interface {
	Resolve(ctx context.Context, expression string) (float64, error)
}
