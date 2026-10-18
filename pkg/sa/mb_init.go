package sa

import (
	"github.com/cloudkucooland/SmartAlac/pkg/mb5"
)

func (c *Curator) initMB5() error {
	if err := mb5.Init(); err != nil {
		return err
	}
	return nil
}
