package service

import (
	"audit-system/ent"
	"audit-system/internal/model"
	"audit-system/internal/utils"
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

type DBService struct {
	client *ent.Client
	mu     sync.Mutex
}

var dbServiceInstance *DBService
var once sync.Once

const defaultDSN = "host=localhost port=5432 user=pq password=pq dbname=audit sslmode=disable"

// GetDBService returns a singleton instance of DBService
func GetDBService() *DBService {
	once.Do(func() {
		dbServiceInstance = &DBService{}
	})
	return dbServiceInstance
}

// Init initializes the database connection and sets up the schema
func (s *DBService) Init(dsn ...string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var usedDSN string
	if len(dsn) > 0 && dsn[0] != "" {
		usedDSN = dsn[0]
	} else {
		usedDSN = defaultDSN
	}

	var err error
	s.client, err = ent.Open("postgres", usedDSN)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}

	// Register the audit log hook
	s.client.Use(AuditLogHook())

	if err := s.client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
}

// Client returns the database client
func (s *DBService) Client() *ent.Client {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.client
}

// Close closes the database connection
func (s *DBService) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.client != nil {
		s.client.Close()
	}
}

func AuditLogHook() ent.Hook {
	return func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, mutation ent.Mutation) (ent.Value, error) {
			auditlogService := GetContainer().AuditLogService
			if ctx.Value(utils.AuditContextKey) == true {
				ctx = context.Background()
				return next.Mutate(ctx, mutation)
			}

			clientID, _ := ctx.Value("clientID").(string)
			entity := mutation.Type()
			op := mutation.Op().String()
			before := map[string]interface{}{}
			after := map[string]interface{}{}

			// Perform the mutation.
			result, err := next.Mutate(ctx, mutation)
			if err != nil {
				return nil, err
			}

			go func() {
				// Create a new context for the audit log operation.
				auditCtx := context.Background()

				// Capture the old value (before the mutation).
				if mutation.Op().Is(ent.OpUpdateOne) {
					fields := mutation.Fields()
					for _, field := range fields {
						value, err := mutation.OldField(ctx, field)
						if err == nil {
							before[field] = value
						}
					}
				}

				// Capture the new value (after the mutation).
				if mutation.Op().Is(ent.OpCreate | ent.OpUpdate | ent.OpUpdateOne | ent.OpDelete | ent.OpDeleteOne) {
					fields := mutation.Fields()
					for _, field := range fields {
						value, exists := mutation.Field(field)
						if exists {
							after[field] = value
						}
					}
				}

				// Create the audit log.
				log := model.AuditLog{
					ClientID:  clientID,
					Timestamp: time.Now(),
					Entity:    entity,
					Mutation:  op,
					Before:    before,
					After:     after,
				}

				auditCtx = context.WithValue(auditCtx, utils.AuditContextKey, true)
				if logErr := auditlogService.CreateAuditLog(auditCtx, log); logErr != nil {
					fmt.Printf("failed to log audit: %v\n", logErr)
				}
			}()

			return result, nil
		})
	}
}
