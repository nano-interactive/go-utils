package testing

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/rzajac/zltest"
)

// Returns new instance of zerolog and *zltest.Tester for testing puurposes
func NewAppTestLogger(t *testing.T, level zerolog.Level) (zerolog.Logger, *zltest.Tester) {
	t.Helper()
	tst := zltest.New(t)
	testWriter := zerolog.NewTestWriter(t)

	logger := zerolog.New(zerolog.MultiLevelWriter(tst, testWriter)).Level(level)

	return logger, tst
}
