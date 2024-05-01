package routes

import (
	"net/http"
	"restapp/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getEventByID(context *gin.Context) {
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse the request"})
		return
	}
	event, err := models.GetEventById(eventId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "No event found for the given ID"})
	}
	context.JSON(http.StatusOK, event)
}

func getEvents(context *gin.Context) {
	events, err := models.GetAllEvents()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not Fetch Events from the Database", "error": err})
		return
	}
	context.JSON(http.StatusOK, events) // sending the response back in JSON format
}

func createEvent(context *gin.Context) {

	var event models.Event
	err := context.ShouldBindJSON(&event)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse the request"})
		return
	}
	event.UserID = context.GetInt64("userId")
	err = event.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not Create Event into the Database"})
		return
	}
	context.JSON(http.StatusCreated, gin.H{"message": "Event Created !", "event": event})
}

func updateEvent(context *gin.Context) {
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse the request"})
		return
	}
	event, err := models.GetEventById(eventId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch event for the given ID"})
		return
	}

	if event.UserID != context.GetInt64("userId") {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "You are not authorized to perform this action"})
		return
	}

	var updatedEvent models.Event

	err = context.ShouldBindJSON(&updatedEvent)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Event values provided by the user"})
		return
	}
	updatedEvent.ID = event.ID
	err = updatedEvent.Update()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not update event for the given ID"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Event Updated !", "event": updatedEvent})
}

func deleteEvent(context *gin.Context) {
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse the request"})
		return
	}
	event, err := models.GetEventById(eventId)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "No Event Found for the given event id"})
		return
	}
	if event.UserID != context.GetInt64("userId") {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "You are not authorized to perform this action"})
		return
	}
	err = models.DeleteEvent(eventId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not Delete event for the given ID"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Event Deleted succesfully!"})
}
