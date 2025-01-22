package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net/smtp"
	"os"

	"github.com/gin-gonic/gin"
)

type MailConfig struct {
	HOST string
	PORT string
	USER string
	PWD  string
}

type ContactForm struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Phone     string `json:"phone" binding:"required"`
	Service   string `json:"service" binding:"required"`
	Message   string `json:"message" binding:"required"`
}

const emailTemplate = `
<p><b>Name:</b> {{.FirstName}} {{.LastName}}</p>
<p><b>Email:</b> {{.Email}}</p>
<p><b>Phone:</b> {{.Phone}}</p>
<p><b>Service Requested:</b> {{.Service}}</p>
<p><b>Message:</b></p>
<p>{{.Message}}</p>
`

func (mc *MailConfig) sendEmail(form ContactForm, recipients []string) error {
	tmpl, err := template.New("email").Parse(emailTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %v", err)
	}

	var emailBody bytes.Buffer
	if err := tmpl.Execute(&emailBody, form); err != nil {
		return fmt.Errorf("failed to execute email template: %v", err)
	}

	toHeader := ""
	for i, recipient := range recipients {
		if i > 0 {
			toHeader += ", "
		}
		toHeader += recipient
	}

	subject := "Contact Form Submission"
	msg := fmt.Sprintf("Subject: %s\nTo: %s\nContent-Type: text/html\n\n%s", subject, toHeader, emailBody.String())

	auth := smtp.PlainAuth("", mc.USER, mc.PWD, mc.HOST)
	err = smtp.SendMail(mc.HOST+":"+mc.PORT, auth, mc.USER, recipients, []byte(msg))
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

func main() {

	mc := &MailConfig{}
	flag.StringVar(&mc.HOST, "MAIL HOST", os.Getenv("EMAIL_HOST"), "MAIL HOST")
	flag.StringVar(&mc.PORT, "MAIL PORT", os.Getenv("EMAIL_PORT"), "MAIL PORT")
	flag.StringVar(&mc.USER, "MAIL USER ", os.Getenv("EMAIL_USER"), "MAIL USER")
	flag.StringVar(&mc.PORT, "MAIL PASSWORD", os.Getenv("EMAIL_PASS"), "MAIL PWD")

	flag.Parse()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.POST("/submit-contact", func(c *gin.Context) {
		var form ContactForm
		if err := c.ShouldBindJSON(&form); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		recipients := []string{
			"lugano.ngulwa@gmail.com",
			"jswigo003@gmail.com",
		}

		if err := mc.sendEmail(form, recipients); err != nil {
			log.Printf("Error sending email: %v", err)
			c.JSON(500, gin.H{"error": "Failed to send email"})
			return
		}

		c.JSON(200, gin.H{"message": "Email sent successfully!"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	slog.Info(fmt.Sprintf("Server is running on port %s", port))
	router.Run(":" + port)
}
