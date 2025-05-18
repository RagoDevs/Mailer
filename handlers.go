package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func (app *application) sendContactEmailHandler(c echo.Context) error {

	var input ContactForm

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, envelope{"error": err.Error()})
	}

	if err := app.validator.Struct(input); err != nil {
		return c.JSON(http.StatusBadRequest, envelope{"error": err.Error()})
	}

	recipients := strings.Split(app.config.mail.recipients, ",")

	if err := app.sendContactUsEmail(input, recipients); err != nil {
		log.Printf("Error sending email: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to send emails"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Emails sent successfully!"})
}

func (app *application) sendWelcomeEmailHandler(c echo.Context) error {

	var input SignupData

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, envelope{"error": err.Error()})
	}

	if err := app.validator.Struct(input); err != nil {
		return c.JSON(http.StatusBadRequest, envelope{"error": err.Error()})
	}

	if err := app.sendWelcomeEmail(input); err != nil {
		log.Printf("Error sending email: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to send email"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Email sent successfully!"})
}

func (app *application) sendActivateEmailHandler(c echo.Context) error {

	var input ActivateOrResetData

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, envelope{"error": err.Error()})
	}

	if err := app.validator.Struct(&input); err != nil {
		return c.JSON(http.StatusBadRequest, envelope{"error": err.Error()})
	}

	if err := app.sendActivateEmail(input); err != nil {
		log.Printf("Error sending email: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to send email"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Email sent successfully!"})
}

func (app *application) sendPasswordResetEmailHandler(c echo.Context) error {

	var input ActivateOrResetData

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, envelope{"error": err.Error()})
	}

	if err := app.validator.Struct(&input); err != nil {
		return c.JSON(http.StatusBadRequest, envelope{"error": err.Error()})
	}

	if err := app.sendPasswordResetEmail(input); err != nil {
		log.Printf("Error sending email: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to send email"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Email sent successfully!"})
}


func (app *application) sendResetCompletedEmailHandler(c echo.Context) error {
	var input ResetCompleteData

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, envelope{"error": err.Error()})
	}

	if err := app.validator.Struct(&input); err != nil {
		return c.JSON(http.StatusBadRequest, envelope{"error": err.Error()})
	}

	if err := app.sendResetCompletedEmail(input); err != nil {
		log.Printf("Error sending email: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to send email"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Email sent successfully!"})
}
