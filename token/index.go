package jwt

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var Jwtkey = []byte("simple_life_gmc_xtt_jwt_key")

type MyClaim struct {
	UserId   int    `json:"id"`
	UserName string `json:"nickname"`
	jwt.StandardClaims
}

func CreateToken(userId int, userName string) (string, error) {
	expiredTIme := time.Now().Add(2 * time.Hour)
	now := time.Now()
	claims := MyClaim{
		UserId:   userId,
		UserName: userName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiredTIme.Unix(),
			IssuedAt:  now.Unix(),
			Issuer:    "gmc-xtt",
		},
	}

	tokenStruct := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenStruct.SignedString(Jwtkey)
}

func CheckToken(token string) (*MyClaim, bool) {
	tokenObj, _ := jwt.ParseWithClaims(token, &MyClaim{}, func(token *jwt.Token) (interface{}, error) {
		return Jwtkey, nil
	})
	if key, _ := tokenObj.Claims.(*MyClaim); tokenObj.Valid {
		return key, true
	} else {
		return nil, false
	}
}

// JwtMiddleware jwt中间件
func JwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		valid := JwtChecker(c)
		if valid {
			c.Next()
		} else {
			c.Abort()
		}
	}
}

func JwtChecker(c *gin.Context) bool {
	tokenStr := c.Request.Header.Get("x-g-token")
	if tokenStr == "" {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "用户不存在"})
		return false
	}
	tokenStruck, ok := CheckToken(tokenStr)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "token不正确"})
		return false
	}
	//token超时
	if time.Now().Unix() > tokenStruck.ExpiresAt {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "token过期"})
		return false
	}
	c.Set("userName", tokenStruck.UserName)
	c.Set("userId", tokenStruck.UserId)
	return true
}
