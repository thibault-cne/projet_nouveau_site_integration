package controllers

import (
	"backend/models"
	calendar_services "backend/services/calendar.services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Return a list of all calendars with the day <= the day parameter
func check_calendars(ctx *gin.Context) {
	day, err := strconv.Atoi(ctx.Param("day"))

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect day"})
		return
	}

	calendars, err := calendar_services.GetAllCalendars(day)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": calendars})
}

func add_calendar(ctx *gin.Context) {
	var calendar models.Calendar

	err := ctx.BindJSON(&calendar)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Internal server error"})
	}

	err = calendar_services.AddCalendar(calendar)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func modifyCalendar(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect id"})
		return
	}

	calendar, err := calendar_services.GetCalendarById(id)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Incorrect id"})
	}

	day := ctx.PostForm("day")
	date := ctx.PostForm("date")
	title := ctx.PostForm("title")
	content := ctx.PostForm("content")

	calendar_services.ModifyCalendar(calendar, date, day, title, content)
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func registerCalendarsRoutes(rg *gin.RouterGroup) {
	router_group := rg.Group("/calendars")
	router_group.GET("/check/:day", check_calendars)
}

func registerAdminCalendarRoutes(rg *gin.RouterGroup) {
	router_group := rg.Group("/calendars")
	router_group.POST("/add/day", add_calendar)
	router_group.POST("/modify/:id", modifyCalendar)
}
