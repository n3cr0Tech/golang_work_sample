package middleware

import (
	"fmt"
	"net/http"
	"time"

	"golang_work_sample/internal/mongodb"
	utils "golang_work_sample/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type MiddlewareHandler struct {
	mongoDB *mongodb.MongoClient
}

func NewMiddlewareHandler(db *mongodb.MongoClient) *MiddlewareHandler {
	return &MiddlewareHandler{mongoDB: db}
}

func (h *MiddlewareHandler) AuthChecker(c *gin.Context) {
	authHeader := c.GetHeader(utils.EnvEntries["AUTH_HEADER"])

	// check auth header
	if len(authHeader) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Auth header missing"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// check token after jwt parsing
	tokenStr := authHeader
	tokenParsed, tokenErr := jwt.Parse(tokenStr, ParseToken)
	if tokenErr != nil || !tokenParsed.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token, try loggin in again for a new one"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// check token claims
	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// check if token has not expired yet
	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	userCollectionsName := utils.EnvEntries["MONGO_USERS_DB"]
	recordIndex := map[string]string{"username": claims["username"].(string)}
	userRecord, userErr := h.mongoDB.FindOne(c.Request.Context(), userCollectionsName, recordIndex)
	if userErr != nil || userRecord == nil {
		errorMsg := fmt.Sprintf("No User record found for: %v", claims["username"])
		c.JSON(http.StatusUnauthorized, gin.H{"error": errorMsg})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("currentUser", userRecord)
	c.Next()
}

func ParseToken(tokenStr *jwt.Token) (interface{}, error) {
	_, tokenIsValid := tokenStr.Method.(*jwt.SigningMethodHMAC)
	if !tokenIsValid {
		return nil, fmt.Errorf("Signing method error: %v", tokenStr.Header["alg"])
	}
	return []byte(utils.EnvEntries["JWT_SECRET"]), nil
}
