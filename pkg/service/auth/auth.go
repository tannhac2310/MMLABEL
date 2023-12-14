package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/casbin/casbin/v2"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"strings"

	"firebase.google.com/go/v4/auth"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/jwtutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/group"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/role"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/user"
)

var (
	ErrNotFoundUsername = errors.New("not found username")
	ErrWrongPassword    = errors.New("wrong password")
)

type Service interface {
	LoginUserNamePassword(ctx context.Context, userName, password string) (*LoginResult, error)
	LoginFirebase(ctx context.Context, idToken string) (*LoginResult, error)
	LoginOTP(ctx context.Context, name, phoneNumber, passcode string) (*LoginResult, error)
	RefreshToken(ctx context.Context, token string) (*LoginResult, error)
	ResetPassword(ctx context.Context, userName, passcode, newPassword string) (*LoginResult, error)
}

type authService struct {
	endforcer    casbin.IEnforcer
	jwtGenerator jwtutil.TokenGenerator
	firebaseAuth *auth.Client

	userService  user.Service
	roleService  role.Service
	groupService group.Service

	userRepo             repository.UserRepo
	userNamePasswordRepo repository.UserNamePasswordRepo
	userFirebaseRepo     repository.UserFirebaseRepo
}

func NewService(
	endforcer casbin.IEnforcer,
	jwtGenerator jwtutil.TokenGenerator,
	firebaseAuth *auth.Client,

	userService user.Service,
	roleService role.Service,
	groupService group.Service,

	userRepo repository.UserRepo,
	userNamePasswordRepo repository.UserNamePasswordRepo,
	userFirebaseRepo repository.UserFirebaseRepo,
) Service {
	return &authService{
		endforcer:    endforcer,
		jwtGenerator: jwtGenerator,
		firebaseAuth: firebaseAuth,

		userService:  userService,
		roleService:  roleService,
		groupService: groupService,

		userRepo:             userRepo,
		userNamePasswordRepo: userNamePasswordRepo,
		userFirebaseRepo:     userFirebaseRepo,
	}
}

func (a *authService) buildCustomClaims(ctx context.Context, userID string) (*jwtutil.Claims, error) {
	roles, err := a.roleService.GetRolesForUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("a.roleService.GetRolesForUser: %w", err)
	}

	groups, err := a.groupService.GetGroupsForUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("a.groupService.GetGroupsForUser: %w", err)
	}

	return &jwtutil.Claims{
		Groups: groups,
		Roles:  roles,
	}, nil
}

func (a *authService) getACL(subjects []string) ([]string, error) {
	records := [][]string{}
	for _, s := range subjects {
		fmt.Println("==============>", s)
		p, err := a.endforcer.GetImplicitPermissionsForUser(s)
		fmt.Println(p)
		if err != nil {
			return nil, fmt.Errorf("r.endforcer.GetImplicitPermissionsForUser: %w", err)
		}

		records = append(records, p...)
	}

	acl := []string{}
	mapCheck := map[string]struct{}{}
	for _, r := range records {
		if len(r) < 2 {
			continue
		}

		p := urlToACL(r[1])
		if p == model.UserRoleRoot {
			acl = []string{p}
			break
		}

		if _, ok := mapCheck[p]; ok {
			continue
		}

		mapCheck[p] = struct{}{}
		acl = append(acl, p)
	}

	return acl, nil
}

func urlToACL(url string) string {
	parts := strings.Split(url, "/")
	if len(parts) != 4 {
		return url
	}

	return strings.Join(parts[2:], "::")
}

type LoginResult struct {
	UserID       string
	Token        string
	RefreshToken string
	ACL          []string
	MainRole     *model.Role
	OtherRoles   []string
}

func (a *authService) buildLoginResult(ctx context.Context, userID string) (*LoginResult, error) {
	claims, err := a.buildCustomClaims(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("a.buildCustomClaims: %w", err)
	}

	jwtToken, refreshToken, err := a.jwtGenerator.Encode(userID, claims)
	if err != nil {
		return nil, fmt.Errorf("err generateJWT: %w", err)
	}

	acl, err := a.getACL(append(claims.Groups, claims.Roles...))
	if err != nil {
		return nil, fmt.Errorf("a.roleService.GetACL: %w", err)
	}

	mainRole, err := a.roleService.HighestRole(ctx, claims.Roles)
	if err != nil && !cockroach.IsErrNoRows(err) {
		return nil, fmt.Errorf("a.roleService.HighestRole: %w", err)
	}

	otherRoles := []string{}
	for _, r := range claims.Roles {
		if r == mainRole.ID {
			continue
		}

		otherRoles = append(otherRoles, r)
	}

	return &LoginResult{
		UserID:       userID,
		Token:        jwtToken,
		RefreshToken: refreshToken,
		ACL:          acl,
		MainRole:     mainRole,
		OtherRoles:   otherRoles,
	}, nil
}
