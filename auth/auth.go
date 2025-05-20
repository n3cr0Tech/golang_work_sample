package auth

import (
	"fmt"
	"net/http"
	"time"

	jwt "example.com/jwt"	
	mongoDB "example.com/mongodb"
	types "example.com/types"
	utils "example.com/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context){
	var req types.Register	
	
	if err := c.BindJSON(&req); err != nil {
		return
	}

	if len(req.Username) == 0 || len(req.Password) == 0{
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Error: username or password field empty"})
		return
	}

	hashedPwd, pwdErr := hashPassword(req.Password)
	if pwdErr != nil{
		c.IndentedJSON(http.StatusOK, gin.H{"message": pwdErr})
		return
	}

	res := createNewUserRecordOnDB(req.Username, hashedPwd)
	if res{
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Created record for user: " + req.Username})
	}else{
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Failed to create record for user: " + req.Username + ". See Server logs"})
	}
}


func Login(c *gin.Context) {
	var req types.Login
	envEntries := utils.EnvEntries

	if err := c.BindJSON(&req); err != nil {
		return
	}

	recordIndex := map[string]string{"username": req.Username}
	userRecord, _recordErr := mongoDB.GetRecord("users", recordIndex)		
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

func createNewUserRecordOnDB(username string, hashedPwd string) bool{
	id := uuid.New()
	newUser := map[string]interface{}{
		"UUID": id.String(),
		"username": username,
		"password": hashedPwd,
		"CreatedAt": time.Now(),
	}
	// var newUser types.User
	// newUser.Username = username
	// newUser.Password = hashedPwd
	// newUser.ID = id.String()
	// newUser.CreatedAt = time.Now()	

	filter := map[string]interface{}{"username": username}
	return mongoDB.EnsureRegisterUser("users", filter, newUser)

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
