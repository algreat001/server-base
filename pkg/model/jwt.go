package model

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"

	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/apiserver/dto"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/config"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/servererrors"
)

type JwtToken struct {
	Token string
	Email string
}

func AuthTokenFromContext(c *gin.Context) *JwtToken {
	token := c.Query("token")
	if token == "" {
		token = strings.Split(c.GetHeader("Authorization")+"  ", " ")[1]
	}
	return &JwtToken{
		Token: token,
		Email: c.Query("email"),
	}
}

func (jt *JwtToken) GetUserFromVerifyAuthToken() (*User, error) {

	token, err := jwt.ParseWithClaims(jt.Token, &dto.TokenClaimsReq{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			err := fmt.Sprintf("unexpected signing method: %v", token.Header["alg"])
			logrus.Warn("invalid token - ", err)
			return nil, errors.New(err)
		}
		return []byte(config.GetInstance().TokenAuthSecurityKey), nil
	})

	if err != nil {
		logrus.Warn("invalid token - ", err, token)
		return nil, servererrors.ErrorInvalidToken
	}

	claims, ok := token.Claims.(*dto.TokenClaimsReq)
	if !ok || !token.Valid {
		logrus.Warn("invalid parse token - ", token)
		return nil, servererrors.ErrorNotAuthenticated
	}

	tokenTime := time.Unix(claims.Expired, 0)

	if tokenTime.Unix() < time.Now().Unix() {
		logrus.Warn("expired user token - ", tokenTime)
		return nil, servererrors.ErrorExpiredAuthToken
	}

	user := NewUserFromJwt(claims, jt.Email)

	return user, nil
}
