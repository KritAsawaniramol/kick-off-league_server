package authentication

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type jwtAuthentication struct {
	secretKey string
}

// fetch user details from the token
// func (j *jwtAuthentication) CurrentUser(c *gin.Context) (*model.User, error) {
// 	err := j.validateJWT(c)
// 	if err != nil {
// 		return &model.User{}, err
// 	}

// 	token, err := j.getToken(c)
// 	if err != nil {
// 		return &model.User{}, err
// 	}

// 	claims, err := token.Claims(jwt.MapClaims)
// 	if err != nil {
// 		return &model.User{}, err
// 	}

// 	userId := uint(claims["user_id"].(float64))

// }

func (j *jwtAuthentication) getToken(c *gin.Context) (*jwt.Token, error) {
	cookie, err := c.Cookie("jwt")
	if err != nil {
		return nil, err
	}
	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})
	return token, err
}

// Auth implements Authentication.
func (j *jwtAuthentication) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := j.getToken(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}
		c.Next()
	}
}

func (j *jwtAuthentication) validateJWT(c *gin.Context) error {

	token, err := j.getToken(c)
	if err != nil {
		return err
	}

	_, ok := token.Claims.(jwt.MapClaims)

	if !token.Valid || !ok {
		return errors.New("invalid token provide")
	}
	return nil
}

// AuthAdmin implements Authentication.
func (j *jwtAuthentication) AuthAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := j.getToken(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
			return
		}

		userRole, ok := claims["role"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
			return
		}

		if err != nil || !token.Valid || userRole != "admin" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		c.Next()
	}
}

// AuthNormalUser implements Authentication.
func (j *jwtAuthentication) AuthNormalUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := j.getToken(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Print("token can't claims")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		userRole, ok := claims["role"]
		if !ok {
			log.Print("not found in token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		if err != nil || !token.Valid || userRole != "normal" {
			log.Print("token not valid")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		userID, ok := claims["user_id"]
		if !ok {
			log.Print("user_id not found in token", claims)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		normalUserID, ok := claims["normal_user_id"]
		if !ok {
			log.Print("normalUser_id not found in token", claims)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		c.Set("user_id", uint(userID.(float64)))
		c.Set("normal_user_id", uint(normalUserID.(float64)))
		c.Next()
	}
}

// AuthOrganizer implements Authentication.
func (j *jwtAuthentication) AuthOrganizer() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := j.getToken(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Println("token claims not ok")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
			return
		}

		userRole, ok := claims["role"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
			return
		}

		if err != nil || !token.Valid || userRole != "organizer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		userID, ok := claims["user_id"]
		if !ok {
			log.Print("user_id not found in token", claims)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		normalUserID, ok := claims["organizer_id"]
		if !ok {
			log.Print("normalUser_id not found in token", claims)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}
		c.Set("user_id", uint(userID.(float64)))
		c.Set("organizer_id", uint(normalUserID.(float64)))
		c.Next()
	}
}

func NewJwtAuthentication(secretKey string) Authentication {
	return &jwtAuthentication{
		secretKey: secretKey,
	}
}

// func authRequired(jwtSecretKey string) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		cookie, err := c.Cookie("jwt")
// 		if err != nil {
// 			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Unauthorized"})
// 			return
// 		}

// 		token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
// 			return []byte(jwtSecretKey), nil
// 		})

// 		if err != nil || !token.Valid {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
// 			return
// 		}
// 		c.Next()
// 	}
// }
