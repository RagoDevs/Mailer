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
	"time"

	"github.com/gin-contrib/cors"
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
<!DOCTYPE html>
<html>
<head>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
        }
        .header {
            background-color: #f8f9fa;
            padding: 20px;
            border-bottom: 2px solid #007bff;
            margin-bottom: 20px;
        }
        .content {
            padding: 0 20px 20px 20px;
        }
        .field {
            margin-bottom: 15px;
        }
        .label {
            font-weight: bold;
            color: #555;
        }
        .message-box {
            background-color: #f8f9fa;
            border-left: 4px solid #007bff;
            padding: 15px;
            margin-top: 15px;
        }
        .footer {
            margin-top: 30px;
            padding-top: 15px;
            border-top: 1px solid #eee;
            font-size: 0.9em;
            color: #777;
        }
    </style>
</head>
<body>
    <div class="header">
        <p>Received on {{.FormattedDate}}</p>
    </div>
    
    <div class="content">
        <div class="field">
            <span class="label">From:</span> {{.FirstName}} {{.LastName}}
        </div>
        
        <div class="field">
            <span class="label">Email:</span> <a href="mailto:{{.Email}}">{{.Email}}</a>
        </div>
        
        <div class="field">
            <span class="label">Phone:</span> {{.Phone}}
        </div>
        
        <div class="field">
            <span class="label">Service Requested:</span> {{.Service}}
        </div>
        
        <div class="field">
            <span class="label">Message:</span>
            <div class="message-box">{{.Message}}</div>
        </div>
        
        <div class="footer">
            <p>This is an automated message from your website contact form.</p>
        </div>
    </div>
</body>
</html>
`

func enableCORS() gin.HandlerFunc {

	corsConfig := cors.DefaultConfig()

	corsConfig.AllowAllOrigins = true

	corsConfig.AllowHeaders = []string{"Content-Type", "Authorization"}

	return cors.New(corsConfig)

}

func (mc *MailConfig) sendEmail(form ContactForm, recipients []string) error {

	type templateData struct {
		ContactForm
		FormattedDate string
	}

	data := templateData{
		ContactForm:   form,
		FormattedDate: time.Now().Format("January 2, 2006 at 3:04 PM"),
	}

	tmpl, err := template.New("email").Parse(emailTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %v", err)
	}

	var emailBody bytes.Buffer
	if err := tmpl.Execute(&emailBody, data); err != nil {
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
	flag.StringVar(&mc.PWD, "MAIL PASSWORD", os.Getenv("EMAIL_PASS"), "MAIL PWD")

	flag.Parse()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(enableCORS())
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
