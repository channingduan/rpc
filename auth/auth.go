package auth

import (
	"context"
	"fmt"
	"github.com/channingduan/rpc/cache"
	"github.com/channingduan/rpc/config"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"strconv"
	"time"
)

type Auth struct {
	config *config.Config
	cache  *cache.Cache
}

func NewAuth(config *config.Config, cache *cache.Cache) *Auth {

	return &Auth{
		config: config,
		cache:  cache,
	}
}

func (a *Auth) CreateToken(id uint) (*Token, error) {

	var token Token
	token.AccessExpire = time.Now().Add(time.Duration(a.config.TokenConfig.AccessExpire) * time.Second).Unix()
	token.RefreshExpire = time.Now().Add(time.Duration(a.config.TokenConfig.RefreshExpire) * time.Second).Unix()
	token.AccessUuid = uuid.New().String()
	token.RefreshUuid = token.AccessUuid + "-" + strconv.Itoa(int(id))
	// 登录认证
	claims := jwt.MapClaims{}
	claims["id"] = id
	claims["authorized"] = true
	claims["access_uuid"] = token.AccessUuid
	claims["expire_time"] = token.AccessExpire
	var err error
	jc := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.AccessToken, err = jc.SignedString([]byte(a.config.TokenConfig.AccessSecret))
	if err != nil {
		return nil, err
	}

	// 刷新认证
	fClaims := jwt.MapClaims{}
	fClaims["id"] = id
	fClaims["refresh_id"] = token.RefreshUuid
	fClaims["refresh_time"] = token.RefreshExpire
	fjc := jwt.NewWithClaims(jwt.SigningMethodHS256, fClaims)
	token.RefreshToken, err = fjc.SignedString([]byte(a.config.TokenConfig.RefreshSecret))
	if err != nil {
		return nil, err
	}

	// 存储写入
	ctx := context.Background()
	if err := a.cache.NewCache().Set(ctx, token.AccessUuid, id, time.Duration(token.AccessExpire)).Err(); err != nil {
		return nil, err
	}
	if err := a.cache.NewCache().Set(ctx, token.RefreshUuid, id, time.Duration(token.RefreshExpire)).Err(); err != nil {
		return nil, err
	}

	return &token, nil
}

func (a *Auth) RefreshToken() {

}

func (a *Auth) ValidToken(token string) error {

	parse, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(a.config.TokenConfig.AccessSecret), nil
	})

	if err != nil {
		return err
	}

	if _, ok := parse.Claims.(jwt.Claims); !ok || !parse.Valid {
		return err
	}

	claims, ok := parse.Claims.(jwt.MapClaims)
	if ok && parse.Valid {
		fmt.Println("claims", claims)
	}
	return nil
}

func (a *Auth) RevokedToken() {

}
