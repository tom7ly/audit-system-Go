package utils

import (
	"audit-system/ent"
	"context"
	"fmt"
)

func WithTx(client *ent.Client, ctx context.Context, fn func(tx *ent.Tx) error) error {
	tx, err := client.Tx(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	if err = fn(tx); err != nil {
		return err
	}

	return nil
}
