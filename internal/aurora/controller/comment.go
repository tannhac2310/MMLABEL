package controller

import (
	"github.com/gin-gonic/gin"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/dto"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/service/comment"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/apperror"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/interceptor"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/routeutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/transportutil"
)

type CommentController interface {
	CreateComment(c *gin.Context)
	EditComment(c *gin.Context)
	DeleteComment(c *gin.Context)
	FindComments(c *gin.Context)
	FindCommentHistories(c *gin.Context)
}

type commentController struct {
	commentService comment.Service
}

func (s *commentController) CreateComment(c *gin.Context) {
	req := &dto.CreateCommentRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	userID := interceptor.UserIDFromCtx(c)

	id, err := s.commentService.CreateComment(c, &comment.CreateCommentOpts{
		UserID:      userID,
		TargetID:    req.TargetID,
		TargetType:  enum.CommentTarget_ProductionPlan, // TODO move it into request
		Content:     req.Content,
		Attachments: []*comment.CreateCommentAttachment{},
	})
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.CreateCommentResponse{
		ID: id,
	})
}

func (s *commentController) EditComment(c *gin.Context) {
	req := &dto.EditCommentRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	userID := interceptor.UserIDFromCtx(c)

	err = s.commentService.EditComment(c, &comment.EditCommentOpts{
		ID:          req.ID,
		UserID:      userID,
		TargetID:    req.TargetID,
		TargetType:  enum.CommentTarget_ProductionPlan, // TODO move it into request
		Content:     req.Content,
		Attachments: []*comment.EditCommentAttachment{},
	})
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.EditCommentResponse{})
}

func (s *commentController) DeleteComment(c *gin.Context) {
	req := &dto.DeleteCommentRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	err = s.commentService.DeleteComment(c, req.ID)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.DeleteCommentResponse{})
}

func (s *commentController) FindComments(c *gin.Context) {
	req := &dto.FindCommentsRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	sort := &repository.Sort{
		Order: repository.SortOrderDESC,
		By:    "ID",
	}
	if req.Sort != nil {
		sort = &repository.Sort{
			Order: repository.SortOrder(req.Sort.Order),
			By:    req.Sort.By,
		}
	}
	comments, cnt, err := s.commentService.FindComments(c, comment.FindCommentsOpts{
		TargetID:   req.Filter.TargetId,
		TargetType: enum.CommentTarget_ProductionPlan,
	}, sort, req.Paging.Limit, req.Paging.Offset)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	commentResp := make([]*dto.Comment, 0, len(comments))
	for _, f := range comments {
		data := &dto.Comment{
			ID:         f.ID,
			UserID:     f.UserID,
			TargetID:   f.TargetID,
			TargetType: f.TargetType,
			Content:    f.Content,
			CreatedAt:  f.CreatedAt,
			UpdatedAt:  f.UpdatedAt,
			DeletedAt:  f.DeletedAt.Time,
			UserName:   f.UserName,
		}

		commentResp = append(commentResp, data)
	}

	transportutil.SendJSONResponse(c, &dto.FindCommentsResponse{
		Comments: commentResp,
		Total:    cnt.Count,
	})
}

func (s *commentController) FindCommentHistories(c *gin.Context) {
	req := &dto.FindCommentHistoriesRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	// sort := &repository.Sort{
	// 	Order: repository.SortOrderDESC,
	// 	By:    "ID",
	// }
	// if req.Sort != nil {
	// 	sort = &repository.Sort{
	// 		Order: repository.SortOrder(req.Sort.Order),
	// 		By:    req.Sort.By,
	// 	}
	// }
	commentHistories, cnt, err := s.commentService.FindCommentHistories(c, req.Filter.ID)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	commentHistoriesResp := make([]*dto.CommentHistory, 0, len(commentHistories))
	for _, f := range commentHistories {
		data := &dto.CommentHistory{
			ID:        f.ID,
			CommentID: f.CommentID,
			Content:   f.Content,
			CreatedAt: f.CreatedAt,
		}

		commentHistoriesResp = append(commentHistoriesResp, data)
	}

	transportutil.SendJSONResponse(c, &dto.FindCommentHistoriesResponse{
		CommentHistories: commentHistoriesResp,
		Total:            cnt.Count,
	})
}

func RegisterCommentController(
	r *gin.RouterGroup,
	commentService comment.Service,
) {
	g := r.Group("comment")

	var c CommentController = &commentController{
		commentService: commentService,
	}

	routeutil.AddEndpoint(
		g,
		"create",
		c.CreateComment,
		&dto.CreateCommentRequest{},
		&dto.CreateCommentResponse{},
		"Create comment",
	)

	routeutil.AddEndpoint(
		g,
		"edit",
		c.EditComment,
		&dto.EditCommentRequest{},
		&dto.EditCommentResponse{},
		"Edit comment",
	)

	routeutil.AddEndpoint(
		g,
		"delete",
		c.DeleteComment,
		&dto.DeleteCommentRequest{},
		&dto.DeleteCommentResponse{},
		"delete comment",
	)

	routeutil.AddEndpoint(
		g,
		"find",
		c.FindComments,
		&dto.FindCommentsRequest{},
		&dto.FindCommentsResponse{},
		"Find comments",
	)

	routeutil.AddEndpoint(
		g,
		"find-history",
		c.FindCommentHistories,
		&dto.FindCommentHistoriesRequest{},
		&dto.FindCommentHistoriesResponse{},
		"Find comment histories",
	)
}
