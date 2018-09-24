package tyrgin

import (
	"time"

	"github.com/appleboy/gin-jwt"
)

// AuthMid returns a pre-configured instance of of the jwt auth middleware
// struct so that we can tell gin servers to use it.
func AuthMid() *jwt.GinJWTMiddleware {
	authMiddleware := &jwt.GinJWTMiddleware{
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

	return authMiddleware
}
