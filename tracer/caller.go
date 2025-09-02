package tracer

import (
	"runtime"
	"strings"
	"sync"
)

var (
	once       sync.Once
	prettifier *PrettyCaller
)

// занимается тем, чтобы красиво определять имя функции, откуда вызван span.
type PrettyCaller struct {
	cutPrefix, cutSuffix string
}

func (c *PrettyCaller) FuncName(skip int) string {
	//достаёт информацию о функции в стеке вызова.
	pc, _, _, _ := runtime.Caller(skip + 1)
	//получает полное имя функции
	name := runtime.FuncForPC(pc).Name()
	if c != nil {
		//отрезает лишние префиксы/суффиксы (например, github.com/user/project/), чтобы имя было более читаемым.
		name = c.pretty(name)
	}
	return name
}

func (c *PrettyCaller) pretty(s string) string {
	if c.cutPrefix != "" {
		s, _ = strings.CutPrefix(s, c.cutPrefix)
	}

	if c.cutSuffix != "" {
		s, _ = strings.CutSuffix(s, c.cutSuffix)
	}
	return s
}

type PrettyCallerCfg struct {
	CutPrefix, CutSuffix string
}

func NewPrettyCaller(cfg *PrettyCallerCfg) *PrettyCaller {
	once.Do(func() {
		prettifier = &PrettyCaller{
			cutPrefix: cfg.CutPrefix,
			cutSuffix: cfg.CutSuffix,
		}
	})
	return prettifier
}
