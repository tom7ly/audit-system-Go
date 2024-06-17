package repository

import (
	"audit-system/ent"
	"context"
	"fmt"
)

func withTx(client *ent.Client, ctx context.Context, fn func(tx *ent.Tx) error) error {
	// Start a transaction
	tx, err := client.Tx(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	// Rollback the transaction in case of error
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

	// Execute the provided function within the transaction
	if err = fn(tx); err != nil {
		return err
	}

	return nil
}
