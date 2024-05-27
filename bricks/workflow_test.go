package bricks_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/janstoon/toolbox/bricks"
)

func TestCompensatorError(t *testing.T) {
	var (
		errOrg = errors.New("original error")
		err1   = errors.New("err one")
		err2   = errors.New("err two")
		err3   = errors.New("err three")

		c      bricks.Compensator
		buffer string
	)

	err := bricks.CompensatorAsError(bricks.CompensatorFunc(func(ctx context.Context, err error) error {
		buffer = "compensated 1"

		return nil
	}), err1)
	require.ErrorAs(t, err, &c)
	require.ErrorIs(t, err, err1)
	require.NoError(t, c.Compensate(context.Background(), errOrg))
	assert.Equal(t, "compensated 1", buffer)

	err = bricks.CompensatorAsError(bricks.CompensatorFunc(func(ctx context.Context, err error) error {
		buffer = "compensated 2"

		return nil
	}), errors.Join(err1, err2))
	require.ErrorAs(t, err, &c)
	require.ErrorIs(t, err, err1)
	require.ErrorIs(t, err, err2)
	require.NoError(t, c.Compensate(context.Background(), errOrg))
	assert.Equal(t, "compensated 2", buffer)

	ece := errors.Join(err1, bricks.CompensatorAsError(bricks.CompensatorFunc(func(ctx context.Context, err error) error {
		buffer = "compensated 3"

		return nil
	}), err2), err3)
	require.ErrorAs(t, ece, &c)
	require.ErrorIs(t, ece, err1)
	require.ErrorIs(t, ece, err2)
	require.ErrorIs(t, ece, err3)
	require.NoError(t, c.Compensate(context.Background(), errOrg))
	assert.Equal(t, "compensated 3", buffer)
}
