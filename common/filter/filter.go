package filter

import (
	"github.com/docker/docker/api/types/filters"
)

func NewFilterArgs(f map[string]string) filters.Args {
	var args []filters.KeyValuePair
	for k, v := range f {
		args = append(args, filters.Arg(k, v))
	}
	return filters.NewArgs(args...)
}
