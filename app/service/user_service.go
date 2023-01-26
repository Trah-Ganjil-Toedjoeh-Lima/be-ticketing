package service

import (
	"errors"
	"github.com/frchandra/gmcgo/app/model"
	"github.com/frchandra/gmcgo/app/repository"
	"github.com/frchandra/gmcgo/app/util"
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

func (us *UserService) GetOrInsertOne(user *model.User) (int64, error) {
	result := us.userRepository.GetOrInsertOne(user)
	if result.Error != nil {
		return 0, nil
	}
	return result.RowsAffected, nil
}

func (us *UserService) InsertOne(user *model.User) (int64, error) {
	//hash the credential //TODO: just for boilerplate. Optional for gmco case
	hashedCred, err := bcrypt.GenerateFromPassword([]byte(user.Phone), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	user.Phone = string(hashedCred)
	//store user to db
	result := us.userRepository.InsertOne(user)
	if result.Error != nil {
		return 0, err
	}
	return result.RowsAffected, nil
}

func (us *UserService) verifyCredentials(cred, hashedCred string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedCred), []byte(cred))
}

func (us *UserService) ValidateLogin(userInput *model.User) error {
	var userOut model.User
	var err error
	//get the user credential pairs email/name & password
	result := us.userRepository.GetByPairs(userInput, &userOut)
	if result.Error != nil {
		return result.Error
	}
	//verify the user credential
	err = us.verifyCredentials(userInput.Phone, userOut.Phone)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return err
	}
	*userInput = userOut
	return nil
}

func (us *UserService) GenerateToken(userInput *model.User) (*util.TokenDetails, error) {
	//create token for us user
	tokenDetails, err := us.tokenUtil.CreateToken(userInput.UserId)
	if err != nil {
		return tokenDetails, err
	}
	//store the token to redis
	if err = us.tokenUtil.StoreAuthn(userInput.UserId, tokenDetails); err != nil {
		return tokenDetails, err
	}
	//return the new created token
	return tokenDetails, nil
}

func (us *UserService) GetById(userId uint64) (model.User, error) {
	var userOut model.User
	result := us.userRepository.GetById(userId, &userOut)
	if result.Error != nil {
		return userOut, result.Error
	}
	if result.RowsAffected < 1 {
		return userOut, errors.New("user tidak ditemukan")
	}
	return userOut, nil
}
