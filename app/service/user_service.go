package service

import (
	"errors"
	"github.com/frchandra/ticketing-gmcgo/app/model"
	"github.com/frchandra/ticketing-gmcgo/app/repository"
	"github.com/frchandra/ticketing-gmcgo/app/util"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepository *repository.UserRepository
	tokenUtil      *util.TokenUtil
}

func NewUserService(userRepository *repository.UserRepository, tokenUtil *util.TokenUtil) *UserService {
	return &UserService{
		tokenUtil:      tokenUtil,
		userRepository: userRepository,
	}
}

func (u *UserService) GetOrInsertOne(user *model.User) (int64, error) {
	result := u.userRepository.GetOrInsertOne(user)
	if result.Error != nil {
		return 0, errors.New("database operation error")
	}
	return result.RowsAffected, nil
}

func (u *UserService) GetByEmail(email string) (model.User, error) {
	var userResult model.User
	result := u.userRepository.GetByEmail(email, &userResult)
	if result.Error != nil {
		return userResult, errors.New("database operation error")
	}
	return userResult, nil
}

func (u *UserService) GetById(userId uint64) (model.User, error) {
	var userResult model.User
	result := u.userRepository.GetById(userId, &userResult)
	if result.Error != nil {
		return userResult, errors.New("database operation error")
	}
	if result.RowsAffected < 1 {
		return userResult, errors.New("cannot find this user")
	}
	return userResult, nil
}

func (u *UserService) InsertOne(user *model.User) (int64, error) {
	result := u.userRepository.InsertOne(user) //store user to db
	if result.Error != nil {
		return 0, errors.New("database operation error")
	}
	return result.RowsAffected, nil
}

func (u *UserService) UpdateById(userId uint64, user *model.User) (int64, error) {
	result := u.userRepository.UpdateById(userId, user)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (u *UserService) verifyCredentials(cred, hashedCred string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedCred), []byte(cred))
}

func (u *UserService) ValidateLogin(userInput *model.User) error {
	var userOut model.User
	var err error

	result := u.userRepository.GetByPairs(userInput, &userOut) //get the user credential pairs email/name & password
	if result.Error != nil {
		return errors.New("database operation error")
	}

	err = u.verifyCredentials(userInput.Phone, userOut.Phone) //verify the user credential
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return errors.New("credential authentication error")
	}
	*userInput = userOut
	return nil
}

func (u *UserService) GenerateToken(userInput *model.User) (*util.TokenDetails, error) {

	tokenDetails, err := u.tokenUtil.CreateToken(userInput.UserId) //create token for the user
	if err != nil {
		return tokenDetails, errors.New("credential authentication error")
	}

	if err = u.tokenUtil.StoreAuthn(userInput.UserId, tokenDetails); err != nil { //store the token to redis
		return tokenDetails, errors.New("credential preparation error")
	}

	return tokenDetails, nil //return the new created token
}
