package states

import "context"

type State interface {
	Run(ctx context.Context)
	Dispose()
}
