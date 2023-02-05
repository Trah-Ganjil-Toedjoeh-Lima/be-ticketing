package util

import (
	"context"
	"errors"
	"fmt"
	"github.com/frchandra/ticketing-gmcgo/config"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"strconv"
	"strings"
	"time"
)

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

type TokenUtil struct {
	db        *redis.Client
	appConfig *config.AppConfig
}

func NewTokenUtil(db *redis.Client, appConfig *config.AppConfig) *TokenUtil {
	return &TokenUtil{
		db:        db,
		appConfig: appConfig,
	}
}

func (tu *TokenUtil) CreateToken(userId uint64) (*TokenDetails, error) {
	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * tu.appConfig.AccessMinute).Unix()
	td.AccessUuid = uuid.New().String()

	td.RtExpires = time.Now().Add(time.Minute * tu.appConfig.RefreshMinute).Unix()
	td.RefreshUuid = uuid.New().String()

	var err error

	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = userId
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(tu.appConfig.AccessSecret))
	if err != nil {
		return nil, err
	}

	//Creating Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = userId
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(tu.appConfig.RefreshSecret))
	if err != nil {
		return nil, err
	}
	return td, nil
}

func (tu *TokenUtil) StoreAuthn(userid uint64, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()
	var ctx = context.Background()

	//store access token (at) to redis
	if err := tu.db.Set(ctx, td.AccessUuid, strconv.Itoa(int(userid)), at.Sub(now)).Err(); err != nil {
		return err
	}
	//store refresh token (rt) to redis
	if err := tu.db.Set(ctx, td.RefreshUuid, strconv.Itoa(int(userid)), rt.Sub(now)).Err(); err != nil {
		return err
	}
	return nil
}

func (tu *TokenUtil) FetchAuthn(uuid string) error {
	var ctx = context.Background()
	//check if token is present in the token storage. get user id from redis given the authentication detail's access uuid
	_, err := tu.db.Get(ctx, uuid).Result()
	if err != nil {
		return err
	}
	return nil
}

func (tu *TokenUtil) DeleteAuthn(givenUuid string) (int64, error) {
	var ctx = context.Background()
	//delete the given uuid from redis
	deleted, err := tu.db.Del(ctx, givenUuid).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}

func (tu *TokenUtil) ExtractToken(c *gin.Context) string {
	//extract token if it on the request param
	token := c.Query("token")
	if token != "" {
		return token
	}
	//extract token if it on the request header
	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

func (tu *TokenUtil) VerifyToken(c *gin.Context, secret string) (*jwt.Token, error) {
	//verify the token format and algorithm
	tokenString := tu.ExtractToken(c)
	if tokenString == "" {
		return nil, errors.New("cannot find token")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok { //Make sure that the token method conform to "SigningMethodHMAC"
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (tu *TokenUtil) ValidateToken(c *gin.Context, secret string) (*jwt.Token, error) {
	//verify the token claims
	token, err := tu.VerifyToken(c, secret)
	if err != nil {
		return token, err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return token, err
	}
	return token, nil
}

type AccessDetails struct {
	AccessUuid string
	UserId     uint64
}

func (tu *TokenUtil) GetValidatedAccess(c *gin.Context) (*AccessDetails, error) {
	token, err := tu.ValidateToken(c, tu.appConfig.AccessSecret)
	if err != nil {
		return nil, err
	}

	claims, _ := token.Claims.(jwt.MapClaims)
	accessUuid, ok := claims["access_uuid"].(string)
	if !ok {
		return nil, err
	}

	userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
	if err != nil {
		return nil, err
	}

	return &AccessDetails{
		AccessUuid: accessUuid,
		UserId:     userId,
	}, nil
}

// TODO: learn tu all

func (tu *TokenUtil) Refresh(c *gin.Context) (map[string]string, error) {
	token, err := tu.ValidateToken(c, tu.appConfig.RefreshSecret)
	if err != nil {
		return nil, err
	}
	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims

	refreshUuid, ok := claims["refresh_uuid"].(string) //convert the interface to string
	if !ok {
		return nil, err
	}
	userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
	if err != nil {
		return nil, err
	}

	//Delete the previous Refresh Token
	deleted, err := tu.DeleteAuthn(refreshUuid)
	if err != nil || deleted == 0 { //if any goes wrong
		return nil, err
	}
	//Create new pairs of refresh and access tokens
	tokenDetails, err := tu.CreateToken(userId)
	if err != nil {
		return nil, err
	}
	//save the tokens metadata to redis
	saveErr := tu.StoreAuthn(userId, tokenDetails)
	if saveErr != nil {
		return nil, err
	}
	tokens := map[string]string{
		"access_token":  tokenDetails.AccessToken,
		"refresh_token": tokenDetails.RefreshToken,
	}
	return tokens, nil
}
