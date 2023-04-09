package auth

import (
	"net/http"
	"regexp"
	"time"

	"github.com/btdjangbah001/chat-app/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("my_secret_key")

type Claims struct {
	UserID uint
	jwt.StandardClaims
}

func createToken(user *models.User) (string, error) {
	// Create JWT claims
	claims := &Claims{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			IssuedAt:  time.Now().Unix(),
			Subject:   "myapp_user",
		},
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign JWT token
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func RegisterUser(c *gin.Context) {
	var user models.User
	var userLogin models.UserSignUp

	err := c.ShouldBindJSON(&userLogin)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userLogin.Password), bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "something went wrong please try again"})
		return
	}

	exists, _ := models.UsernameOrEmailExists(&userLogin)
	if exists {
		c.JSON(400, gin.H{"error": "username or email already exists"})
		return
	}

	user.Username = userLogin.Username
	user.Email = userLogin.Email
	user.Password = string(hashedPassword)

	err = user.CreateUser()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	token, err := createToken(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "something went wrong please try again"})
		return
	}

	c.JSON(200, gin.H{"user": user, "token": token})
}

func LoginUser(c *gin.Context) {
	var userLogin models.UserLogin
	var user *models.User
	var err error

	err = c.ShouldBindJSON(&userLogin)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if userFieldIsEmail(userLogin.UserField) {
		user, err = models.GetUserByEmail(userLogin.UserField)
		if err != nil {
			// try to get user by username if not found by email
			user, err = models.GetUserByUsername(userLogin.UserField)
			if err != nil {
				c.JSON(400, gin.H{"error": "invalid email or username or password"})
				return
			}
		}
	} else {
		user, err = models.GetUserByUsername(userLogin.UserField)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid email or username or password"})
			return
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userLogin.Password))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid email or username or password"})
		return
	}

	token, err := createToken(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "something went wrong please try again"})
		return
	}

	c.JSON(200, gin.H{"user": user, "token": token})
}

func userFieldIsEmail(emailOrUsername string) bool {
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return re.MatchString(emailOrUsername)
}
