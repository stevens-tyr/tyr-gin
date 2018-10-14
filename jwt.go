package tyrgin

import (
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

func Authenticator(c *gin.Context) (interface{}, error) {
	var login Login
	if err := c.ShouldBindJSON(&login); err != nil {
		return "Missing login values.", jwt.ErrMissingLoginValues
	}

	s, err := GetSession()
	if err != nil {
		return "Can not get Mongo Session.", MongoSessionFailure
	}

	db := GetDataBase("tyr", s)
	col, err := SafeGetCollection("users", db)
	if err != nil {
		return "Mongo Collection does not exist.", MongoCollectionFailure
	}

	var user User
	if err = col.Find(bson.M{"email": login.Email}).One(&user); err != nil {
		return "User not found.", UserNotFoundError
	}

	if err = bcrypt.CompareHashAndPassword(user.Password, []byte(login.Password)); err != nil {
		return "Incorrect password", IncorrectPasswordError
	}

	return "Success", nil
}

func Authorizator(d interface{}, c *gin.Context) bool {
	var email Email
	if err := c.ShouldBindJSON(&email); err != nil {
		return false
	}

	s, err := GetSession()
	if err != nil {
		return false
	}

	db := GetDataBase("tyr", s)
	col, err := SafeGetCollection("users", db)
	if err != nil {
		return false
	}

	var user User
	if err = col.Find(bson.M{"email": email.Email}).One(&user); err != nil {
		return false
	}

	return true
}

func Register(c *gin.Context) {
	var register Register
	if err := c.ShouldBindJSON(&register); err != nil {
		return
	}

	s, err := GetSession()
	if err != nil {
		return
	}

	db := GetDataBase("tyr", s)
	col, err := SafeGetCollection("users", db)
	if err != nil {
		return
	}

	var user User
	if err = col.Find(bson.M{"email": register.Email}).One(&user); err != nil {
		return
	}
}

func Unauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"status_code": code,
		"message":     message,
	})
}

var authMiddleware = &jwt.GinJWTMiddleware{
	Realm:          "localhost:5555",
	Key:            []byte("tyr makes you tear"),
	Timeout:        time.Hour,
	MaxRefresh:     time.Hour * 24,
	Authenticator:  Authenticator,
	Authorizator:   Authorizator,
	Unauthorized:   Unauthorized,
	TokenLookup:    "header:Authorization",
	TokenHeadName:  "Bearer",
	TimeFunc:       time.Now,
	SendCookie:     true,
	SecureCookie:   false,
	CookieHTTPOnly: true,
	CookieDomain:   "localhost:5555",
	CookieName:     "token",
	TokenLookup:    "cookie:token",
}

// AuthMid returns a pre-configured instance of of the jwt auth middleware
// struct so that we can tell gin servers to use it.
func AuthMid() *jwt.GinJWTMiddleware {

	return authMiddleware
}
