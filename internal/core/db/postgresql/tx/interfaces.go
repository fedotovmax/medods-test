package tx

import (
	"context"

	"github.com/fedotovmax/medods-test/internal/core/db/postgresql"
)

type Extractor interface {
	ExtractTx(ctx context.Context) postgresql.Executor
}
