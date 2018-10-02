package tyrgin

import (
	"time"

	jwt "github.com/appleboy/gin-jwt"
)

var authMiddleware = &jwt.GinJWTMiddleware{
	Realm:      "localhost:3000",
	Key:        []byte("secret key"),
	Timeout:    time.Hour,
	MaxRefresh: time.Hour * 24,
	//Authenticator: Authenticator,
	//Authorizator:  Authorizator,
	//Unauthorized:  Unauthorized,
	TokenLookup:   "header:Authorization",
	TokenHeadName: "Bearer",
	TimeFunc:      time.Now,
}

// AuthMid returns a pre-configured instance of of the jwt auth middleware
// struct so that we can tell gin servers to use it.
func AuthMid() *jwt.GinJWTMiddleware {

	return authMiddleware
}
