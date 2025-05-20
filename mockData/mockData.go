package mockData

import (
	"fmt"
	"net/http"

	types "example.com/types"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"rsc.io/quote"
)

// album represents data about a record album.
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// albums slice to seed record album data.
var albums = []album{
	{ID: "0", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "1", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "2", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

// getAlbums responds with the list of all albums as JSON.
func GetAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

func findAlbum(albums []album, id string) (album, string) {
	for _, e := range albums {
		if id == e.ID {
			return e, ""
		}
	}
	var empty album
	errMsg := fmt.Sprintf("No album found for id: %v", id)
	return empty, errMsg
}

func GetAlbum(c *gin.Context) {
	id := c.Param("id")

	var egress album
	egress, err := findAlbum(albums, id)

	if len(err) == 0 {
		c.IndentedJSON(http.StatusOK, egress)
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"message": err})
	}

}

// ----- Auth Placeholders -----
var mockUsers []types.User

func CreateMockUsers() {
	user1 := types.User{UUID: "0", Username: "mickeyMouse"}
	user2 := types.User{UUID: "1", Username: "foo-user"}

	Mock_Pwd := "password123"
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(Mock_Pwd), bcrypt.DefaultCost)

	if err != nil {
		fmt.Println(err)
		return
	}
	user1.Password = string(passwordHash)
	user2.Password = string(passwordHash)
	mockUsers = append(mockUsers, user1, user2)
}

func GetUserByUsername(username string) (*types.User, string) {
	for _, u := range mockUsers {
		if username == u.Username {
			return &u, ""
		}
	}	
	errMsg := fmt.Sprintf("No record found for user: %v", username)
	return nil, errMsg
}

// postAlbums adds an album from JSON received in the request body.
func PostGreeting(c *gin.Context) {
	var reqPayload types.QuoteIngress

	// Call BindJSON to bind the received JSON to
	// resPayload
	if err := c.BindJSON(&reqPayload); err != nil {
		return
	}

	var resPayload types.QuoteEgress
	resPayload.Message = generateGreeting(reqPayload.Name)
	c.IndentedJSON(http.StatusCreated, resPayload)
}

func generateGreeting(name string) string {
	msg := fmt.Sprintf("Hello %v, here's a a quote: %v", name, quote.Glass())
	return msg
}

// ---------------

/*
--- SAMPLE REQUEST ---
curl http://localhost:8080/albums

curl http://localhost:8080/albums/2
*/
