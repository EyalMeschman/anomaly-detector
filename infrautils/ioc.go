package infrautils

import (
	"strings"

	"go.uber.org/dig"
)

func IocProvideWrapper(c *dig.Container, constructor any) {
	if err := c.Provide(constructor); err != nil {
		if !strings.Contains(err.Error(), "already provided by") {
			panic(err)
		}
	}
}
