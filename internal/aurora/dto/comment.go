package dto

import (
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/commondto"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

type CommentFilter struct {
	TargetId string `json:"targetId,omitempty"`
}

type FindCommentsRequest struct {
	Filter *CommentFilter    `json:"filter" binding:"required"`
	Paging *commondto.Paging `json:"paging" binding:"required"`
	Sort   *commondto.Sort   `json:"sort"`
}

type FindCommentsResponse struct {
	Comments []*Comment `json:"comments"`
	Total    int64      `json:"total"`
}

type CommentHistoryFilter struct {
	ID string `json:"id,omitempty"`
}

type FindCommentHistoriesRequest struct {
	Filter *CommentHistoryFilter `json:"filter" binding:"required"`
	Paging *commondto.Paging     `json:"paging" binding:"required"`
	Sort   *commondto.Sort       `json:"sort"`
}

type FindCommentHistoriesResponse struct {
	CommentHistories []*CommentHistory `json:"commentHistories,omitempty"`
	Total            int64             `json:"total,omitempty"`
}

type Attachment struct {
	FileURL string `json:"fileURL,omitempty" binding:"required"`
}

type CreateCommentRequest struct {
	TargetID    string             `json:"targetID,omitempty"`
	TargetType  enum.CommentTarget `json:"targetType,omitempty"`
	Content     string             `json:"content,omitempty" binding:"required"`
	Attachments []*Attachment      `json:"attachments,omitempty"`
}

type CreateCommentResponse struct {
	ID string `json:"id"`
}

type EditCommentRequest struct {
	ID          string        `json:"id"`
	TargetID    string        `json:"targetID,omitempty"`
	Content     string        `json:"content,omitempty" binding:"required"`
	Attachments []*Attachment `json:"attachments,omitempty"`
}

type EditCommentResponse struct{}

type DeleteCommentRequest struct {
	ID string `json:"id"`
}

type DeleteCommentResponse struct{}

type Comment struct {
	ID         string    `json:"id,omitempty"`
	UserID     string    `json:"userID,omitempty"`
	TargetID   string    `json:"targetID,omitempty"`
	TargetType int16     `json:"targetType,omitempty"`
	Content    string    `json:"content,omitempty"`
	CreatedAt  time.Time `json:"createdAt,omitempty"`
	UpdatedAt  time.Time `json:"updatedAt,omitempty"`
	DeletedAt  time.Time `json:"deletedAt,omitempty"`
	UserName   string    `json:"userName,omitempty"`
}

type CommentHistory struct {
	ID        string    `json:"id,omitempty"`
	CommentID string    `json:"commentID,omitempty"`
	Content   string    `json:"content,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
}
