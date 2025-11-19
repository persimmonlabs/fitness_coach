package domain

import (
	"time"

	"github.com/google/uuid"
)

// Conversation represents a conversation thread with the AI coach
type Conversation struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index:idx_user_conversations" json:"user_id"`
	Title     *string    `gorm:"type:varchar(255)" json:"title,omitempty"` // Auto-generated or user-set
	Context   *string    `gorm:"type:jsonb" json:"context,omitempty"` // JSON context for the conversation

	CreatedAt time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP;index:idx_user_conversations" json:"created_at"`
	UpdatedAt time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	User     User      `gorm:"foreignKey:UserID" json:"-"`
	Messages []Message `gorm:"foreignKey:ConversationID;constraint:OnDelete:CASCADE" json:"messages,omitempty"`
}

// TableName specifies the table name for GORM
func (Conversation) TableName() string {
	return "conversations"
}

// Message represents a message in a conversation
type Message struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ConversationID uuid.UUID `gorm:"type:uuid;not null;index:idx_conversation_messages" json:"conversation_id"`
	Role           string    `gorm:"type:varchar(50);not null" json:"role"` // user, assistant, system
	Content        string    `gorm:"type:text;not null" json:"content"`
	Metadata       *string   `gorm:"type:jsonb" json:"metadata,omitempty"` // Additional metadata (e.g., tokens, model)

	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;index:idx_conversation_messages" json:"created_at"`

	// Relationships
	Conversation Conversation `gorm:"foreignKey:ConversationID" json:"-"`
}

// TableName specifies the table name for GORM
func (Message) TableName() string {
	return "messages"
}
