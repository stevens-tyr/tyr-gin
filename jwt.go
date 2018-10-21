package tyrgin

import (
	"regexp"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

// Authenticator a default function for a gin jwt, that authenticates a user.
func Authenticator(c *gin.Context) (interface{}, error) {
	var login Login
	if err := c.ShouldBindJSON(&login); err != nil {
		return "Missing login values.", jwt.ErrMissingLoginValues
	}

	s, err := GetSession()
	if err != nil {
		return "Can not get Mongo Session.", ErrorMongoSessionFailure
	}

	db := GetDataBase("tyr", s)
	col, err := SafeGetCollection("users", db)
	if err != nil {
		return "Mongo Collection does not exist.", ErrorMongoCollectionFailure
	}

	var user User
	if err = col.Find(bson.M{"email": login.Email}).One(&user); err != nil {
		return "User not found.", ErrorUserNotFound
	}

	if err = bcrypt.CompareHashAndPassword(user.Password, []byte(login.Password)); err != nil {
		return "Incorrect password", ErrorIncorrectPassword
	}

	return "Success", nil
}

// Authorizator a default function for a gin jwt, that authorizes a user.
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

// Register a function that registers a User.
func Register(c *gin.Context) {
	var register RegisterForm
	if err := c.ShouldBindJSON(&register); err != nil {
		return
	}

	s, err := GetSession()
	if err != nil {
		c.JSON(500, gin.H{
			"status_code": 500,
			"message":     "Failed to get Mongo Session.",
		})
		return
	}

	db := GetDataBase("tyr", s)
	col, err := SafeGetCollection("users", db)
	if err != nil {
		c.JSON(500, gin.H{
			"status_code": 500,
			"message":     "Failed to get collection.",
		})
		return
	}

	if err = IsValidEmail(register.Email); err != nil {
		c.JSON(400, gin.H{
			"status_code": 400,
			"message":     "Email is invalid.",
		})
		return
	}

	var user User
	if err = col.Find(bson.M{"email": register.Email}).One(&user); err != nil {
		c.JSON(400, gin.H{
			"status_code": 400,
			"message":     "Email is taken.",
		})
		return
	}

	if register.Password != register.PasswordConfirmation {
		c.JSON(400, gin.H{
			"status_code": 400,
			"message":     "Your password and password confirmation do not match.",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{
			"status_code": 500,
			"message":     "Failed to generate hash",
		})
		return
	}

	user = User{
		Email:    register.Email,
		Password: hash,
		First:    register.First,
		Last:     register.Last,
		Roles:    make([]string, 0),
	}

	err = col.Insert(&user)
	if err != nil {
		c.JSON(500, gin.H{
			"status_code": 500,
			"message":     "Failed to create user.",
		})
		return
	}

	c.JSON(200, gin.H{
		"status_code": 200,
		"message":     "User created.",
	})

}

// Unauthorized a default jwt gin function, called when authentication is failed.
func Unauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"status_code": code,
		"message":     message,
	})
}

var authMiddleware = &jwt.GinJWTMiddleware{
	Realm:         "localhost:5555",
	Key:           []byte("tyr makes you tear"),
	Timeout:       time.Hour,
	MaxRefresh:    time.Hour * 24,
	Authenticator: Authenticator,
	Authorizator:  Authorizator,
	Unauthorized:  Unauthorized,
	TokenLookup:   "header:Authorization",
	TokenHeadName: "Bearer",
	TimeFunc:      time.Now,
	SendCookie:    true,
	SecureCookie:  false,
	//CookieHTTPOnly: true,
	//CookieDomain:   "localhost:5555",
	//CookieName:     "token",
}

// AuthMid returns a pre-configured instance of of the jwt auth middleware
// struct so that we can tell gin servers to use it.
func AuthMid() *jwt.GinJWTMiddleware {

	return authMiddleware
}

// IsValidEmail checks an email string to be valid using regex
// returns ErrorEmailNotValid
// TODO: add host check or send validation email
func IsValidEmail(email string) error {
	validEmailForm := regexp.MustCompile("\\A[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\z")
	if !validEmailForm.MatchString(email) {
		return ErrorEmailNotValid
	}
	return nil

}
