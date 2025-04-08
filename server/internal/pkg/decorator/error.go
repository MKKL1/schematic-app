package decorator

import (
	"context"
	"fmt"
)

type commandErrorDecorator[C any, R any] struct {
	base CommandHandler[C, R]
}

func (d commandErrorDecorator[C, R]) Handle(ctx context.Context, cmd C) (res R, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("executing command %s: %w", generateActionName(cmd), err)
		}
	}()

	return d.base.Handle(ctx, cmd)
}

type queryErrorDecorator[C any, R any] struct {
	base QueryHandler[C, R]
}

func (d queryErrorDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("executing query %s: %w", generateActionName(cmd), err)
		}
	}()

	return d.base.Handle(ctx, cmd)
}
