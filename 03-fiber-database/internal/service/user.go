package service

import (
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

func (us UserService) GetUserByID(username string) (domain.User, error) {
	eu, err := us.userRepo.GetOneByID(username)
	if err != nil {
		return domain.User{}, err
	}

	return domain.NewUserFromEntity(eu), nil
}

// REFACTOR:
func (us UserService) Login(username string, password string) (TokenString, error) {
	user, err := us.userRepo.GetOneByID(username)
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

	return us.authSrv.GenerateToken(user.Username)
}

func (us UserService) Register(username string, password string) error {
	user, err := us.userRepo.GetOneByID(username)
	if err != nil && err != repository.ErrUserNotFound {
		return err
	}

	if !user.IsEmpty() {
		return ErrUserAlreadyRegistered
	}

	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	if err != nil {
		return err
	}

	return us.authSrv.userRepo.Insert(username, string(hashedPasswordBytes))
}

func (us UserService) List() ([]domain.User, error) {
	eusers, err := us.userRepo.ListAll()
	if err != nil {
		return []domain.User{}, err
	}

	dusers := make([]domain.User, len(eusers))
	for idx, eu := range eusers {
		dusers[idx] = domain.NewUserFromEntity(eu)
	}

	return dusers, nil
}
