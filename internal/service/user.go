package service

import (
	"context"
	"errors"

	"example.com/authorization/internal/domain"
	"example.com/authorization/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo repository.UserRepository
	authSrv  AuthService
}

func NewUserService(userRepo repository.UserRepository, authSrv AuthService) UserService {
	return UserService{
		userRepo: userRepo,
		authSrv:  authSrv,
	}
}

func (us UserService) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
	eu, err := us.userRepo.GetOneByUsername(ctx, username)
	if err != nil {
		return domain.User{}, err
	}

	return domain.NewUserFromEntity(eu), nil
}

func (us UserService) GetUserByID(ctx context.Context, UserID int64) (domain.User, error) {
	eu, err := us.userRepo.GetOneByID(ctx, UserID)
	if err != nil {
		return domain.User{}, err
	}

	return domain.NewUserFromEntity(eu), nil
}

func (us UserService) Login(ctx context.Context, username string, password string) (TokenString, error) {
	user, err := us.userRepo.GetOneByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return "", errors.Join(repository.ErrUserNotFound, ErrUserNotFound)
		}
		return "", err
	}

	cerr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if cerr != nil {
		return "", errors.Join(ErrWrongCredentials, cerr)
	}

	return us.authSrv.GenerateToken(user.Id)
}

func (us UserService) Register(ctx context.Context, username string, password string) error {
	user, err := us.userRepo.GetOneByUsername(ctx, username)
	if err != nil && !errors.Is(repository.ErrUserNotFound, err) {
		return err
	}

	if user.IsValid() {
		return ErrUserAlreadyRegistered
	}

	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	if err != nil {
		return err
	}

	return us.authSrv.userRepo.Insert(ctx, username, string(hashedPasswordBytes))
}

func (us UserService) List(ctx context.Context) ([]domain.User, error) {
	eusers, err := us.userRepo.ListAll(ctx)
	if err != nil {
		return []domain.User{}, err
	}

	dusers := make([]domain.User, len(eusers))
	for idx, eu := range eusers {
		dusers[idx] = domain.NewUserFromEntity(eu)
	}

	return dusers, nil
}
