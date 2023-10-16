package testing

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/rzajac/zltest"
)

func NewAppTestLogger(t testing.TB, level zerolog.Level) (zerolog.Logger, *zltest.Tester) {
	t.Helper()
	tst := zltest.New(t)
	testWriter := zerolog.NewTestWriter(t)
	logger := zerolog.New(zerolog.MultiLevelWriter(tst, testWriter)).Level(level)

	return logger, tst
}
