package auth

import (
	"net/http"
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
	var userRegister models.UserSignUp

	err := c.ShouldBindJSON(&userRegister)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !models.UserFieldIsEmail(userRegister.Email) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "enter a valid email"})
		return
	}

	if userRegister.Password != userRegister.ConfirmPassword {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "passwords do not match"})
		return
	}

	// if !isStrongPassword(userLogin.Password) {
	// 	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "password is not strong enough"})
	// 	return
	// }

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRegister.Password), bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "something went wrong please try again"})
		return
	}

	exists, _ := models.UsernameOrEmailExists(&userRegister)
	if exists {
		c.JSON(400, gin.H{"error": "username or email already exists"})
		return
	}

	user.Username = userRegister.Username
	user.Email = userRegister.Email
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

	if models.UserFieldIsEmail(userLogin.UserField) {
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

/*
This regex matches a password string that:

is 8 to 15 characters long
contains at least one lowercase letter (a-z)
contains at least one uppercase letter (A-Z)
contains at least one digit (0-9)
contains at least one special character that is not alphanumeric (e.g., !@#$%^&*)
*/
// func isStrongPassword(password string) bool {
// 	re := regexp.MustCompile(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[^\da-zA-Z]).{8,15}$`)
// 	return re.MatchString(password)
// }
