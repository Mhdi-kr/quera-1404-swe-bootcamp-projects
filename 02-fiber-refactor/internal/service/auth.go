package service

import "example.com/authorization/internal/repository"

type AuthService struct {
	jwtSecret string
	userRepo  repository.UserRepository
}

func NewAuthorizationService(jwtSecret string, userRepo repository.UserRepository) AuthService {
	return AuthService{
		jwtSecret: jwtSecret,
		userRepo:  userRepo,
	}
}

func (as AuthService) Validate() error {
	return nil
}

// TODO: move this to service layer
// hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), 4)
// if err != nil {
// 	return err
// }

// func checkUser(user entity.User) bool {
// 	isUserFound := false

// 	for _, u := range usersDB {
// 		resCompare := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password))
// 		fmt.Println(resCompare)
// 		if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password)) == nil && user.Username == u.Username {
// 			isUserFound = true
// 		}
// 	}

// 	return isUserFound
// }
