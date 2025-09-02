package tracer

import (
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCaller(t *testing.T) {
	const (
		prefix = "gitlab.telespace/"
	)

	t.Run("test singleton", func(t *testing.T) {
		_ = NewPrettyCaller(&PrettyCallerCfg{
			CutPrefix: prefix,
		})

		p := NewPrettyCaller(nil)
		assert.Equal(t, prefix, p.cutPrefix)
	})

	t.Run("test cutting", func(t *testing.T) {
		p := NewPrettyCaller(&PrettyCallerCfg{CutPrefix: prefix})

		pc, _, _, _ := runtime.Caller(0)
		funcName := runtime.FuncForPC(pc).Name()
		expected, _ := strings.CutPrefix(funcName, prefix)
		assert.Equal(t, expected, p.FuncName(0))
	})

	t.Run("nil caller", func(t *testing.T) {
		pc, _, _, _ := runtime.Caller(0)
		expectedFuncName := runtime.FuncForPC(pc).Name()

		assert.Equal(t, expectedFuncName, prettifier.FuncName(0))
	})
}
