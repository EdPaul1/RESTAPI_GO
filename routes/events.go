package routes

import (
	"log"
	"net/http"
	"strconv"

	"example.com/rest-api/models"
	"github.com/gin-gonic/gin"
)

func getEvent(context *gin.Context) {
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid event id"})
		return
	}

	event, err := models.GetEventByID(eventId)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch event"})
		return
	}
	context.JSON(http.StatusOK, event)
}
func getEvents(context *gin.Context) {
	events, err := models.GetAllEvents()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, events)
}

func createEvent(context *gin.Context) {

	var event models.Event
	err := context.ShouldBindJSON(&event)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Could not parse request data."})
		return
	}
	userId := context.GetInt64("userId")
	event.UserID = userId

	err = event.Save()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create event. Try Again Later"})
		return
	}
	context.JSON(http.StatusCreated, gin.H{"message": "Event Created!", "event": event})
}

func updateEvent(context *gin.Context) {
	log.Println("Received request to update an event.")

	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		log.Printf("Error parsing event ID: %v", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse event id."})
		return
	}
	log.Printf("Parsed event ID: %d", eventId)

	userId := context.GetInt64("userId")

	event, err := models.GetEventByID(eventId)

	if err != nil {
		log.Printf("Error fetching event with ID %d: %v", eventId, err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch the event."})
		return
	}

	if event.UserID != userId {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "You are not authorized to update this event."})
		return
	}

	var updatedEvent models.Event
	log.Printf("Type of updatedEvent: %T", updatedEvent) // Debugging the type
	err = context.ShouldBindJSON(&updatedEvent)          // Ensure pointer is passed
	if err != nil {
		log.Printf("Error parsing request body: %v", err)
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse request data.",
			"error":   err.Error(),
		})
		return
	}
	log.Printf("Parsed request body successfully: %+v", updatedEvent)

	updatedEvent.ID = eventId
	err = updatedEvent.Update()
	if err != nil {
		log.Printf("Error updating event with ID %d: %v", eventId, err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not update event."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Event updated successfully!"})
	log.Println("Response sent: Event updated successfully.")
}

func deleteEvent(context *gin.Context) {
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid event id"})
		return
	}

	userId := context.GetInt64("userId")
	event, err := models.GetEventByID(eventId)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch event"})
		return
	}

	if event.UserID != userId {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "You are not authorized to delete this event."})
		return
	}

	err = event.Delete()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete event"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully!"})

}
