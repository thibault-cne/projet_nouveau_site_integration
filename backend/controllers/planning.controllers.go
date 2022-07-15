package controllers

import (
	planning_services "backend/services/planning.services"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func addPlanning(ctx *gin.Context) {
	// Get the begining and end date from the frontend
	beginingDate := ctx.PostForm("beginingDate")
	endDate := ctx.PostForm("endDate")

	// Get the file from the frontend
	file, err := ctx.FormFile("file")

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file"})
		return
	}

	// Get the file extension
	fileExtension := filepath.Ext(file.Filename)

	// Check if the file extension is valid
	if !(fileExtension == ".jpg" || fileExtension == ".jpeg" || fileExtension == ".png") {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file extension"})
		return
	}

	// Get the file size
	fileSize := file.Size

	// Check if the file size is valid
	if fileSize > 5000000 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "File too big"})
		return
	}

	// Get the file content
	fileContent, err := file.Open()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error while opening the file"})
		return
	}

	// Create a new file
	pictureName := "planning_pictures_" + strings.Replace(beginingDate, "/", "-", -1) + strings.Replace(endDate, "/", "-", -1) + fileExtension
	newFile, err := os.Create("static/images/planning_pictures/" + pictureName)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot create the file"})
		return
	}

	// Copy the file content to the new file
	_, err = io.Copy(newFile, fileContent)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error while copying the file"})
		return
	}

	// Close the file
	err = newFile.Close()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error while closing the file"})
		return
	}

	// Add the planning to the database
	err = planning_services.AddPlaning(planning_services.NewPlaning(pictureName, beginingDate, endDate))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error while adding the planning"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func retriveCurrentPlanning(ctx *gin.Context) {
	currentPlanning, err := planning_services.RetrieveLastPlanning()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// planningResponse, err := api_response_services.NewPlanningResponse(currentPlanning)

	fileContent, err := ioutil.ReadFile("static/images/planning_pictures/" + currentPlanning.Picture)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.Data(http.StatusOK, "image/jpeg", fileContent)
}

func RegisterPlanningRoutes(rg *gin.RouterGroup) {
	route_group := rg.Group("/planning")
	route_group.POST("/add", addPlanning)
	route_group.GET("/current", retriveCurrentPlanning)
}