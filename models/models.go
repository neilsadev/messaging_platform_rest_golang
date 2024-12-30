package models

import (
	"time"
)

// User Model
type User struct {
	UserUUID       string `gorm:"primaryKey"`
	DepartmentUUID string `gorm:"type:uuid"`
	BranchUUID     string `gorm:"type:uuid"`
	Email          string `gorm:"unique;not null"`
	Phone          string `gorm:"not null"`
	Username       string `gorm:"unique;not null"`
	AvatarURL      string
	CreatedAt      time.Time  `gorm:"autoCreateTime"`
	GroupsCreated  []Group    `gorm:"foreignKey:CreatedBy"`
	MessagesSent   []Message  `gorm:"foreignKey:SenderUUID"`
	Reactions      []Reaction `gorm:"foreignKey:UserUUID"`
}

// Group Model
type Group struct {
	GroupUUID   string `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Description string
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	CreatedBy   string    `gorm:"type:uuid"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	UpdatedBy   string    `gorm:"type:uuid"`
	DeletedAt   *time.Time
	DeletedBy   string            `gorm:"type:uuid"`
	Memberships []GroupMembership `gorm:"foreignKey:GroupUUID"`
	Messages    []Message         `gorm:"foreignKey:GroupUUID"`
}

// GroupMembership Model
type GroupMembership struct {
	MembershipUUID string    `gorm:"primaryKey"`
	GroupUUID      string    `gorm:"type:uuid;not null"`
	UserUUID       string    `gorm:"type:uuid;not null"`
	JoinedAt       time.Time `gorm:"autoCreateTime"`
	LeftAt         *time.Time
}

// Message Model
type Message struct {
	MessageUUID   string    `gorm:"primaryKey"`
	SenderUUID    string    `gorm:"type:uuid;not null"`
	RecipientUUID *string   `gorm:"type:uuid"`
	GroupUUID     *string   `gorm:"type:uuid"`
	Content       string    `gorm:"type:text;not null"`
	Attachments   string    `gorm:"type:json"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
	DeletedAt     *time.Time
	Reactions     []Reaction      `gorm:"foreignKey:MessageUUID"`
	Threads       []MessageThread `gorm:"foreignKey:ParentMessageUUID"`
}

// MessageThread Model
type MessageThread struct {
	ThreadUUID        string    `gorm:"primaryKey"`
	ParentMessageUUID string    `gorm:"type:uuid;not null"`
	ChildMessageUUID  string    `gorm:"type:uuid;not null"`
	CreatedAt         time.Time `gorm:"autoCreateTime"`
}

// Reaction Model
type Reaction struct {
	ReactionUUID string    `gorm:"primaryKey"`
	MessageUUID  string    `gorm:"type:uuid;not null"`
	UserUUID     string    `gorm:"type:uuid;not null"`
	Reaction     string    `gorm:"type:string;not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}
