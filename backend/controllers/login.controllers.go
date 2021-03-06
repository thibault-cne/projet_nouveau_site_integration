package controllers

import (
	"backend/models"
	claims_services "backend/services/claims.services"
	users_services "backend/services/users.services"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func validate_Google_OAuth_token(ctx *gin.Context) {
	access_token := ctx.PostForm("credential")

	resp, err := http.Get("https://oauth2.googleapis.com/tokeninfo?id_token=" + access_token)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error validating token"})
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var google_api_response models.GoogleApiResponse
	json.Unmarshal(body, &google_api_response)

	login_api_response, err_message, err := create_response(google_api_response.Email, google_api_response.Name)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err_message})
	}

	ctx.JSON(http.StatusOK, login_api_response)
}

func create_response(email string, name string) (*models.LoginApiResponse, string, error) {
	user, err := users_services.GetUserByEmail(email)

	if err != nil {
		user_id, err := users_services.AddUser(users_services.NewUser(email, name))

		if err != nil {
			return nil, "Internal server error", err
		}

		return claims_services.NewLoginApiResponse(user_id), "", nil
	} else {
		user_id := user.ID

		return claims_services.NewLoginApiResponse(user_id), "", nil
	}
}

func refresh_token(ctx *gin.Context) {
	reqToken := ctx.Request.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	reqToken = splitToken[1]

	new_access_token, err := claims_services.Refresh_token(reqToken)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
	}

	ctx.JSON(http.StatusOK, gin.H{"access_token": new_access_token})
}

func registerLoginRoutes(rg *gin.RouterGroup) {
	router_group := rg.Group("/login")
	router_group.POST("/g-oauth", validate_Google_OAuth_token)
	router_group.GET("/refresh-token", refresh_token)
}
