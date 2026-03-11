package auth

import (
	"fmt"
	"net/http"
	"time"

	jwt "golang_work_sample/internal/jwt"
	mongoDB "golang_work_sample/internal/mongodb"
	types "golang_work_sample/internal/types"
	utils "golang_work_sample/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	mongoDB *mongoDB.MongoClient
}

func NewAuthHandler(db *mongoDB.MongoClient) *AuthHandler {
	return &AuthHandler{mongoDB: db}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req types.Register

	if err := c.BindJSON(&req); err != nil {
		return
	}

	if len(req.Username) == 0 || len(req.Password) == 0 {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Error: username or password field empty"})
		return
	}

	hashedPwd, pwdErr := hashPassword(req.Password)
	if pwdErr != nil {
		c.IndentedJSON(http.StatusOK, gin.H{"message": pwdErr})
		return
	}

	res := h.createNewUserRecordOnDB(req.Username, hashedPwd)
	if res {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Created record for user: " + req.Username})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Failed to create record for user: " + req.Username + ". See Server logs"})
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req types.Login
	envEntries := utils.EnvEntries

	if err := c.BindJSON(&req); err != nil {
		return
	}

	recordIndex := map[string]string{"username": req.Username}
	userRecord, _recordErr := h.mongoDB.GetRecord("users", recordIndex)
	if _recordErr != nil {
		fmt.Println("INSIDE - record error clause")
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Failed to find record for " + req.Username})
		return
	}

	pwdErr := bcrypt.CompareHashAndPassword([]byte(userRecord.Password), []byte(req.Password))
	if pwdErr != nil {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Invalid password"})
		return
	}

	jwtToken, _err := jwt.CreateToken(envEntries["JWT_SECRET"], userRecord.Username)
	fmt.Println("jwt token: ", jwtToken)
	if _err == nil {
		c.IndentedJSON(http.StatusOK, gin.H{"token: ": jwtToken})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"message": _err})
	}
}

func (h *AuthHandler) createNewUserRecordOnDB(username string, hashedPwd string) bool {
	id := uuid.New()
	newUser := map[string]interface{}{
		"UUID":      id.String(),
		"username":  username,
		"password":  hashedPwd,
		"CreatedAt": time.Now(),
	}

	filter := map[string]interface{}{"username": username}
	return h.mongoDB.EnsureRegisterUser("users", filter, newUser)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

/*
--- SAMPLE REQUEST ---
curl http://localhost:8080/login \
    --include --header \
    "Content-Type: application/json" \
    --request "POST" --data \
    '{"username": "foo-user", "password": "password123"}'
*/
