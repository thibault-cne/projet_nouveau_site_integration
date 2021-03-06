package controllers

import (
	"net/http"

	claims_services "backend/services/claims.services"
	parrainageservices "backend/services/parrainage.services"
	usersservices "backend/services/users.services"

	"github.com/gin-gonic/gin"
)

func ensureLoggedIn() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		loggedInInterface, _ := ctx.Get("is_logged_in")
		loggedIn := loggedInInterface.(bool)
		if !loggedIn {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token."})
		}
	}
}

func setUserStatus() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reqToken := ctx.Request.Header.Get("Authorization")
		claims, err := claims_services.RetrieveUserClaims(reqToken)

		if err != nil {
			ctx.Set("is_logged_in", false)
			return
		}

		ctx.Set("is_logged_in", true)
		ctx.Set("user_id", claims.User_id)
	}
}

func ensureUserIsAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userIdInterface, _ := ctx.Get("user_id")
		userId := userIdInterface.(int)

		userIsAdmin, _ := usersservices.CheckAdmin(userId)

		if !userIsAdmin {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "You are not administrator."})
		}
	}
}

func ensureParrainageProcessIsOn(parrainageProcess *parrainageservices.ParrainageProcess) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !parrainageProcess.IsProcessOpen {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Parrainage process is closed."})
		}
	}
}

func ensureUserNotInYearOne() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userIdInterface, _ := ctx.Get("user_id")
		userId := userIdInterface.(int)

		user, _ := usersservices.GetUser(userId)

		if user.Current_year == 1 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "You are not allowed to access parrainage process."})
		}
	}
}
