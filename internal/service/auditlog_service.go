package service

import (
	"audit-system/internal/model"
	"audit-system/internal/repository"
	"context"
	"time"
)

type AuditLogService struct {
	repo *repository.AuditLogRepository
}

func NewAuditLogService(repo *repository.AuditLogRepository) *AuditLogService {
	return &AuditLogService{repo: repo}
}

func (s *AuditLogService) CreateAuditLog(ctx context.Context, log model.AuditLog) error {
	return s.repo.CreateAuditLog(ctx, log)
}

func (s *AuditLogService) GetAllAuditLogs(ctx context.Context) ([]*model.AuditLog, error) {
	return s.repo.GetAllAuditLogs(ctx)
}

func (s *AuditLogService) GetAuditLogsByEmail(ctx context.Context, email string) ([]*model.AuditLog, error) {
	return s.repo.GetAuditLogsByEmail(ctx, email)
}

func (s *AuditLogService) DeleteOldAuditLogs(ctx context.Context, ttl time.Duration) (int, error) {
	return s.repo.DeleteOldAuditLogs(ctx, ttl)
}
