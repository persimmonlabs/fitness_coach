package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"fitness-tracker/internal/core/domain"
	"fitness-tracker/internal/core/ports"
)

type conversationRepository struct {
	db *gorm.DB
}

// NewConversationRepository creates a new conversation repository
func NewConversationRepository(db *gorm.DB) ports.ConversationRepository {
	return &conversationRepository{db: db}
}

func (r *conversationRepository) Create(ctx context.Context, conversation *domain.Conversation) error {
	return r.db.WithContext(ctx).Create(conversation).Error
}

func (r *conversationRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Conversation, error) {
	var conversation domain.Conversation
	err := r.db.WithContext(ctx).
		Preload("Messages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		Where("id = ?", id).
		First(&conversation).Error
	if err != nil {
		return nil, err
	}
	return &conversation, nil
}

func (r *conversationRepository) Update(ctx context.Context, conversation *domain.Conversation) error {
	return r.db.WithContext(ctx).Save(conversation).Error
}

func (r *conversationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// Cascade delete is handled by the database constraint
	return r.db.WithContext(ctx).Delete(&domain.Conversation{}, "id = ?", id).Error
}

func (r *conversationRepository) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domain.Conversation, error) {
	var conversations []*domain.Conversation
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Limit(limit).
		Offset(offset).
		Order("updated_at DESC").
		Find(&conversations).Error
	if err != nil {
		return nil, err
	}
	return conversations, nil
}

// Message operations

func (r *conversationRepository) AddMessage(ctx context.Context, message *domain.Message) error {
	return r.db.WithContext(ctx).Create(message).Error
}

func (r *conversationRepository) GetMessages(ctx context.Context, conversationID uuid.UUID, limit, offset int) ([]*domain.Message, error) {
	var messages []*domain.Message
	err := r.db.WithContext(ctx).
		Where("conversation_id = ?", conversationID).
		Limit(limit).
		Offset(offset).
		Order("created_at ASC").
		Find(&messages).Error
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *conversationRepository) GetLatestMessages(ctx context.Context, conversationID uuid.UUID, limit int) ([]*domain.Message, error) {
	var messages []*domain.Message
	err := r.db.WithContext(ctx).
		Where("conversation_id = ?", conversationID).
		Limit(limit).
		Order("created_at DESC").
		Find(&messages).Error
	if err != nil {
		return nil, err
	}

	// Reverse the slice to get chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}
