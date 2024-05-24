package controllers

import (
	"encoding/json"
	"fmt"
	"go-backend/models"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type Event struct {
	ClientId string `json:"clientId"`
	Type     string `json:"type"`
	Source   string `json:"source"`
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) GetEventbyPK(context fiber.Ctx) error {
	id := context.Params("id")
	eventModel := &models.Events{}

	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "No ID provided.",
		})
		return nil
	}

	fmt.Println("EVENT_ID => ", id)

	err := r.DB.Where("id = ?", id).First(eventModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Could not fetch event by id provided",
		})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Event fetched successfully",
		"data":    eventModel,
	})

	return nil
}

func (r *Repository) CreateEvent(context fiber.Ctx) error {
	event := Event{}
	bodyBytes := context.Body()

	if err := json.Unmarshal(bodyBytes, &event); err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "Invalid request body"})
		return err
	}

	err := r.DB.Create(&event).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "Failed persisting Event on database."})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "Event created successfully."})

	return nil
}

func (r *Repository) GetEvents(context fiber.Ctx) error {
	eventModels := &[]models.Events{}

	err := r.DB.Find(eventModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "Failed on fetching events."})
		return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message": "All events fetched successfully.",
			"data":    eventModels,
		})
	return nil
}

func (r *Repository) DeleteEvent(context fiber.Ctx) error {
	eventModel := models.Events{}
	id := context.Params("id")

	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "No ID provided for deletion.",
		})
		return nil
	}

	err := r.DB.Delete(eventModel, id)

	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Event could not be deleted.",
		})
		return err.Error
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Event deleted successfully.",
	})

	return nil
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Get("/getEvent/:id", r.GetEventbyPK)
	api.Get("/getEvents", r.GetEvents)
	api.Post("/createEvent", r.CreateEvent)
	api.Delete("deleteEvent/:id", r.DeleteEvent)
}
