package runner

import (
	"github.com/mlctrez/godom/nap"
)

func Run(o *nap.Config) (err error) {
	if err = setup(o); err != nil {
		o.Logger.Debug("setup error", "error", err)
		return err
	}
	return run(o)
}
