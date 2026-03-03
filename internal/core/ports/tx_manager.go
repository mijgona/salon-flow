package ports

import "context"

// TxManager provides transactional boundaries.
type TxManager interface {
	Execute(ctx context.Context, fn func(tx interface{}) error) error
}
