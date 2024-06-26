package hooks

// import (
// 	"audit-system/ent"
// 	"audit-system/internal/model"
// 	"audit-system/internal/utils"
// 	"context"
// 	"fmt"
// 	"time"
// )

// // AuditLogHook is a hook for logging mutations.
// func AuditLogHook() ent.Hook {
// 	return func(next ent.Mutator) ent.Mutator {
// 		return ent.MutateFunc(func(ctx context.Context, mutation ent.Mutation) (ent.Value, error) {
// 			if ctx.Value(utils.AuditContextKey) == true {
// 				ctx = context.Background()
// 				return next.Mutate(ctx, mutation)
// 			}

// 			clientID, _ := ctx.Value("clientID").(string)
// 			entity := mutation.Type()
// 			op := mutation.Op().String()
// 			before := map[string]interface{}{}
// 			after := map[string]interface{}{}

// 			// Perform the mutation.
// 			result, err := next.Mutate(ctx, mutation)
// 			if err != nil {
// 				return nil, err
// 			}

// 			go func() {
// 				// Create a new context for the audit log operation.
// 				auditCtx := context.Background()

// 				// Capture the old value (before the mutation).
// 				if mutation.Op().Is(ent.OpUpdateOne) {
// 					fields := mutation.Fields()
// 					for _, field := range fields {
// 						value, err := mutation.OldField(ctx, field)
// 						if err == nil {
// 							before[field] = value
// 						}
// 					}
// 				}

// 				// Capture the new value (after the mutation).
// 				if mutation.Op().Is(ent.OpCreate | ent.OpUpdate | ent.OpUpdateOne | ent.OpDelete | ent.OpDeleteOne) {
// 					fields := mutation.Fields()
// 					for _, field := range fields {
// 						value, exists := mutation.Field(field)
// 						if exists {
// 							after[field] = value
// 						}
// 					}
// 				}

// 				// Create the audit log.
// 				log := model.AuditLog{
// 					ClientID:  clientID,
// 					Timestamp: time.Now(),
// 					Entity:    entity,
// 					Mutation:  op,
// 					Before:    before,
// 					After:     after,
// 				}

// 				auditCtx = context.WithValue(auditCtx, utils.AuditContextKey, true)
// 				if logErr := auditLogService.CreateAuditLog(auditCtx, log); logErr != nil {
// 					fmt.Printf("failed to log audit: %v\n", logErr)
// 				}
// 			}()

// 			return result, nil
// 		})
// 	}
// }
