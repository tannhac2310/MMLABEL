package user

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"

	"firebase.google.com/go/v4/auth"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/group"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/role"
)

type Service interface {
	FindUserByID(ctx context.Context, id string) (*model.User, error)
	CreateFirebaseUser(ctx context.Context, firebaseID string, opts *CreateUserOpts) (userID string, err error)
	CreateUser(ctx context.Context, opts *CreateUserOpts) (userID string, err error)
	EditUserProfile(ctx context.Context, opts *EditUserProfileOpts) error
	CheckExistedUserName(ctx context.Context, userName string) (bool, error)
	CreateLoginAccount(ctx context.Context, userID, email, phoneNumber, password string) error
	ChangePassword(ctx context.Context, userID, password string) error
	SearchUsers(ctx context.Context, opts *SearchUsersOpts, limit, offset int64) ([]*repository.UserData, *repository.CountResult, error)
	VerifyPassword(ctx context.Context, userID, password string) (bool, error)
	UpdateFCMToken(ctx context.Context, userID, deviceID, token string) error
	SoftDelete(ctx context.Context, id string) error
}

type userService struct {
	firebaseAuth *auth.Client
	roleService  role.Service
	groupService group.Service

	userRepo             repository.UserRepo
	userRoleRepo         repository.UserRoleRepo
	userNamePasswordRepo repository.UserNamePasswordRepo
	userGroupRepo        repository.UserGroupRepo
	userFirebaseRepo     repository.UserFirebaseRepo
	userFCMTokenRepo     repository.UserFCMTokenRepo
}

func NewService(
	firebaseAuth *auth.Client,
	roleService role.Service,
	groupService group.Service,
	userRepo repository.UserRepo,
	userRoleRepo repository.UserRoleRepo,
	userNamePasswordRepo repository.UserNamePasswordRepo,
	userGroupRepo repository.UserGroupRepo,
	userFirebaseRepo repository.UserFirebaseRepo,
	userFCMTokenRepo repository.UserFCMTokenRepo,
) Service {
	return &userService{
		firebaseAuth:         firebaseAuth,
		roleService:          roleService,
		groupService:         groupService,
		userRepo:             userRepo,
		userRoleRepo:         userRoleRepo,
		userNamePasswordRepo: userNamePasswordRepo,
		userGroupRepo:        userGroupRepo,
		userFirebaseRepo:     userFirebaseRepo,
		userFCMTokenRepo:     userFCMTokenRepo,
	}
}
