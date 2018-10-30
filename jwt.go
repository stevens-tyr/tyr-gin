package tyrgin

import (
	"fmt"
	"os"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/goware/emailx"
	godotenv "github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	bcrypt "golang.org/x/crypto/bcrypt"
	mgo "gopkg.in/mgo.v2"
	bson "gopkg.in/mgo.v2/bson"
)

var authMiddleware *jwt.GinJWTMiddleware

func init() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}

	if env == "dev" {
		if err := godotenv.Load(".dev.env"); err != nil {
			fmt.Println("wtf")
			log.Fatal("Could not load .dev.env.")
		}
	} else {
		if prodErr := godotenv.Load(); prodErr != nil {
			fmt.Println("wtf")
			log.Fatal("Could not load .env.")
		}
	}

	authMiddleware = &jwt.GinJWTMiddleware{
		Realm:         os.Getenv("JWT_REALM"),
		Key:           []byte(os.Getenv("JWT_SECRET")),
		Timeout:       time.Hour,
		MaxRefresh:    time.Hour * 24,
		Authenticator: Authenticator,
		Authorizator:  Authorizator,
		PayloadFunc:   PayloadFunc,
		Unauthorized:  Unauthorized,
		TokenLookup:   "header:Authorization",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
		SendCookie:    true,
		SecureCookie:  false,
		//CookieHTTPOnly: true,
		//CookieDomain: "localhost:5555",
		//CookieName:   "token",
	}
}

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

	db := GetDataBase(os.Getenv("DB_NAME"), s)
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

	return user, nil
}

// Authorizator a default function for a gin jwt, that authorizes a user.
func Authorizator(d interface{}, c *gin.Context) bool {
	//claims := jwt.ExtractClaims(c)

	// todo in future check
	// look at route see if what part of course/assignment they are accessing
	// check if they have permission claims["courses"] slice of enrolledCourse

	return true
}

// Register a function that registers a User.
func Register(c *gin.Context) {
	var register RegisterForm
	if err := c.ShouldBindJSON(&register); err != nil {
		c.JSON(400, gin.H{
			"status_code": 400,
			"message":     "Incorrect json format.",
		})
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

	db := GetDataBase(os.Getenv("DB_NAME"), s)
	col, err := SafeGetCollection("users", db)
	if err != nil {
		c.JSON(500, gin.H{
			"status_code": 500,
			"message":     "Failed to get collection.",
		})
		return
	}

	if err = IsValidEmail(register.Email); err != nil {
		msg := "Email is invalid"
		if err == ErrorUnresolvableEmailHost {
			msg = "Unable to resolve email host"
		}
		c.JSON(400, gin.H{
			"status_code": 400,
			"message":     msg,
		})
		return
	}

	var user User
	if err = col.Find(bson.M{"email": register.Email}).One(&user); err != mgo.ErrNotFound {
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
		Email:           register.Email,
		Password:        hash,
		First:           register.First,
		Last:            register.Last,
		EnrolledCourses: make([]enrolledCourse, 0),
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

// PayloadFunc uses the User's courses as jwt claims.
func PayloadFunc(data interface{}) jwt.MapClaims {
	switch data.(type) {
	case User:
		return jwt.MapClaims{"courses": data.(User).EnrolledCourses}
	default:
		return jwt.MapClaims{}
	}
}

// Unauthorized a default jwt gin function, called when authentication is failed.
func Unauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"status_code": code,
		"message":     message,
	})
}

// AuthMid returns a pre-configured instance of of the jwt auth middleware
// struct so that we can tell gin servers to use it.
func AuthMid() *jwt.GinJWTMiddleware {

	return authMiddleware
}

// IsValidEmail checks an email string to be valid and with resolvable host
func IsValidEmail(email string) error {

	err := emailx.Validate(email)
	if err != nil {
		if err == emailx.ErrInvalidFormat {
			return ErrorEmailNotValid
		}
		if err == emailx.ErrUnresolvableHost {
			return ErrorUnresolvableEmailHost
		}
		return err
	}
	return nil
}
