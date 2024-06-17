package repository

import (
	"audit-system/ent"
	"audit-system/ent/auditlog"
	"audit-system/internal/model"
	"context"
)

type AuditLogRepository struct {
	client *ent.Client
}

func NewAuditLogRepository(client *ent.Client) *AuditLogRepository {
	return &AuditLogRepository{client: client}
}

func (r *AuditLogRepository) CreateAuditLog(ctx context.Context, log model.AuditLog) error {
	_, err := r.client.AuditLog.Create().
		SetClientID(log.ClientID).
		SetTimestamp(log.Timestamp).
		SetEntity(log.Entity).
		SetMutation(log.Mutation).
		SetBefore(log.Before).
		SetAfter(log.After).
		Save(ctx)
	return err
}

func (r *AuditLogRepository) GetAllAuditLogs(ctx context.Context) ([]*model.AuditLog, error) {
	logs, err := r.client.AuditLog.Query().All(ctx)
	if err != nil {
		return nil, err
	}

	var result []*model.AuditLog
	for _, l := range logs {
		result = append(result, &model.AuditLog{
			ClientID:  l.ClientID,
			Timestamp: l.Timestamp,
			Entity:    l.Entity,
			Mutation:  l.Mutation,
			Before:    l.Before,
			After:     l.After,
		})
	}
	return result, nil
}

func (r *AuditLogRepository) GetAuditLogsByEmail(ctx context.Context, email string) ([]*model.AuditLog, error) {
	logs, err := r.client.AuditLog.Query().Where(auditlog.ClientID(email)).All(ctx)
	if err != nil {
		return nil, err
	}

	var result []*model.AuditLog
	for _, l := range logs {
		result = append(result, &model.AuditLog{
			ClientID:  l.ClientID,
			Timestamp: l.Timestamp,
			Entity:    l.Entity,
			Mutation:  l.Mutation,
			Before:    l.Before,
			After:     l.After,
		})
	}
	return result, nil
}
