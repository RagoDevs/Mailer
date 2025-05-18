package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type MailConfig struct {
	HOST       string
	PORT       string
	USER       string
	PWD        string
	RECEPIENTS string
	ALLOWED_IP string
}

type ContactForm struct {
	FirstName string `json:"first_name" form:"first_name" query:"first_name" param:"first_name" validate:"required"`
	LastName  string `json:"last_name" form:"last_name" query:"last_name" param:"last_name" validate:"required"`
	Email     string `json:"email" form:"email" query:"email" param:"email" validate:"required,email"`
	Phone     string `json:"phone" form:"phone" query:"phone" param:"phone" validate:"required"`
	Service   string `json:"service" form:"service" query:"service" param:"service" validate:"required"`
	Message   string `json:"message" form:"message" query:"message" param:"message" validate:"required"`
}

type SignupData struct {
	ID    string `json:"id" form:"id" query:"id" param:"id" validate:"required"`
	Email string `json:"email" form:"email" query:"email" param:"email" validate:"required,email"`
	Token string `json:"token" form:"token" query:"token" param:"token" validate:"required"`
}

type ActivateOrResetData struct {
	Email string `json:"email" form:"email" query:"email" param:"email" validate:"required,email"`
	Token string `json:"token" form:"token" query:"token" param:"token" validate:"required"`
}

type ResetCompleteData struct {
	Email string `json:"email" form:"email" query:"email" param:"email" validate:"required,email"`
}

// CustomValidator is a validator implementation for Echo using go-playground/validator
type CustomValidator struct {
	validator *validator.Validate
}

// Validate validates the provided struct
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

// validateRequest validates the request data
func validateRequest(c echo.Context, data interface{}) error {
	if err := c.Bind(data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(data); err != nil {
		return err
	}

	return nil
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

const welcome_template = `
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
            text-align: center;
        }
        .logo {
            max-width: 200px;
            margin-bottom: 10px;
        }
        .content {
            padding: 0 20px 20px 20px;
        }
        .button {
            background-color: #007bff;
            color: white !important;    
            padding: 12px 24px;
            text-decoration: none;
            border-radius: 4px;
            font-weight: bold;
            display: inline-block;
            margin: 20px 0;
            text-shadow: none;         
            opacity: 1;            
        }
        .button:hover {
            background-color: #0069d9;
        }
        .important-note {
            background-color: #f8f9fa;
            border-left: 4px solid #ffc107;
            padding: 15px;
            margin: 15px 0;
            font-size: 0.9em;
        }
        .footer {
            margin-top: 30px;
            padding-top: 15px;
            border-top: 1px solid #eee;
            font-size: 0.9em;
            color: #777;
            text-align: center;
        }
        .contact-info {
            margin-top: 20px;
        }
    </style>
</head>
<body>
    <div class="header">
        <h2>Welcome to Rent Management System</h2>
    </div>
    
    <div class="content">
  
        <p>Thank you for choosing Rent Management System for your property management needs. We're delighted to welcome you to our platform.</p>
        
        <p>Your account has been created successfully with the following details:</p>
        <p><strong>User ID:</strong> {{.ID}}</p>
        
        <h3>Important: Please Activate Your Account</h3>
        
        <p>To complete your registration and access all features of our platform, please activate your account by clicking the button below:</p>
        
        <a href="https://rent.ragodevs.com/activate?token={{.Token}}" class="button">Activate Account</a>
        
        <div class="important-note">
            <p>Please note that this activation link will expire in 3 days and can only be used once.</p>
        </div>
        
        <div class="contact-info">
            <p>If you have any questions or need assistance, please contact our support team:</p>
            <p>Email: <a href="mailto:support@ragodevs.com">support@ragodevs.com</a><br>
            Phone: +255 654 051 622</p>
        </div>
        
        <p>We look forward to helping you streamline your property management operations.</p>
        
        <p>Best regards,</p>
        <p>The Rent Management System Team</p>
        
        <div class="footer">
            <p><a href="https://rent.ragodevs.com">rent.ragodevs.com</a></p>
            <p>© 2025 Rent Management System. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`

const activate_template = `
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
            text-align: center;
        }
        .logo {
            max-width: 200px;
            margin-bottom: 10px;
        }
        .content {
            padding: 0 20px 20px 20px;
        }

        .button {
            background-color: #007bff;
            color: white !important;    
            padding: 12px 24px;
            text-decoration: none;
            border-radius: 4px;
            font-weight: bold;
            display: inline-block;
            margin: 20px 0;
            text-shadow: none;         
            opacity: 1;            
        }

        .button:hover {
            background-color: #0069d9;
        }
        .important-note {
            background-color: #f8f9fa;
            border-left: 4px solid #ffc107;
            padding: 15px;
            margin: 15px 0;
            font-size: 0.9em;
        }
        .footer {
            margin-top: 30px;
            padding-top: 15px;
            border-top: 1px solid #eee;
            font-size: 0.9em;
            color: #777;
            text-align: center;
        }
        .contact-info {
            margin-top: 20px;
        }
    </style>
</head>
<body>
    <div class="header">
        <h2>Account Activation Required</h2>
    </div>
    
    <div class="content">
        
        <h3>Important: Please Activate Your Account</h3>
        
        <p>Please activate your account to continuing using our services by clicking the button below:</p>
        
        <a href="https://rent.ragodevs.com/activate?token={{.Token}}" class="button">Activate Account</a>
        
        <div class="important-note">
            <p>Please note that this activation link will expire in 3 days and can only be used once.</p>
        </div>
        
        <div class="contact-info">
            <p>If you have any questions or need assistance, please contact our support team:</p>
            <p>Email: <a href="mailto:support@ragodevs.com">support@ragodevs.com</a><br>
            Phone: +255 654 051 622</p>
        </div>
        
        <p>We look forward to helping you streamline your property management operations.</p>
        
        <p>Best regards,</p>
        <p>The Rent Management System Team</p>
        
        <div class="footer">
            <p><a href="https://rent.ragodevs.com">www.ragodevs.com</a></p>
            <p>© 2025 Rent Management System. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`

const pwdreset_template = `
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
            text-align: center;
        }
        .logo {
            max-width: 200px;
            margin-bottom: 10px;
        }
        .content {
            padding: 0 20px 20px 20px;
        }

       .button {
         background-color: #007bff;
         color: white !important;    
         padding: 12px 24px;
         text-decoration: none;
         border-radius: 4px;
         font-weight: bold;
         display: inline-block;
         margin: 20px 0;
         text-shadow: none;         
         opacity: 1;            
        }

        .button:hover {
            background-color: #0069d9;
        }
        .important-note {
            background-color: #f8f9fa;
            border-left: 4px solid #ffc107;
            padding: 15px;
            margin: 15px 0;
            font-size: 0.9em;
        }
        .footer {
            margin-top: 30px;
            padding-top: 15px;
            border-top: 1px solid #eee;
            font-size: 0.9em;
            color: #777;
            text-align: center;
        }
        .contact-info {
            margin-top: 20px;
        }
    </style>
</head>
<body>
    <div class="header">
        <h2>Password Reset Request</h2>
    </div>
    
    <div class="content">
        
        
        <p>Your password can be reset by clicking the button below. If you did not request a new password, please ignore this email.</p>
        
        <a href="https://rent.ragodevs.com/reset?token={{.Token}}" class="button">Reset Password</a>
        
        <div class="important-note">
            <p>Please note that this reset link will expire 45 minutes and can only be used once.</p>
        </div>
        
        <div class="contact-info">
            <p>If you have any questions or need assistance, please contact our support team:</p>
            <p>Email: <a href="mailto:support@ragodevs.com">support@ragodevs.com</a><br>
            Phone: +255 654 051 622</p>
        </div>
        
        <p>Best regards,</p>
        <p>The Rent Management System Team</p>
        
        <div class="footer">
            <p><a href="https://rent.ragodevs.com">rent.ragodevs.com</a></p>
            <p>© 2025 Rent Management System. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`

const completedreset_template = `
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
            text-align: center;
        }
        .logo {
            max-width: 200px;
            margin-bottom: 10px;
        }
        .content {
            padding: 0 20px 20px 20px;
        }

        .button {
            background-color: #007bff;
            color: white !important;    
            padding: 12px 24px;
            text-decoration: none;
            border-radius: 4px;
            font-weight: bold;
            display: inline-block;
            margin: 20px 0;
            text-shadow: none;         
            opacity: 1;            
        }
        .button:hover {
            background-color: #0069d9;
        }
        .important-note {
            background-color: #f8f9fa;
            border-left: 4px solid #ffc107;
            padding: 15px;
            margin: 15px 0;
            font-size: 0.9em;
        }
        .footer {
            margin-top: 30px;
            padding-top: 15px;
            border-top: 1px solid #eee;
            font-size: 0.9em;
            color: #777;
            text-align: center;
        }
        .contact-info {
            margin-top: 20px;
        }
    </style>
</head>
<body>
    <div class="header">
        <h2>Password Changed</h2>
    </div>
    
    <div class="content">
        
        
        <p>You have successfully changed your password for Rent Management System.</p>
        <p>If this wasn't done by you, please immediately reset the password of your Rent Management System.</p>

        <div class="contact-info">
            <p>If you have any questions or need assistance, please contact our support team:</p>
            <p>Email: <a href="mailto:support@ragodevs.com">support@ragodevs.com</a><br>
            Phone: +255 654 051 622</p>
        </div>
        
        <p>Best regards,</p>
        <p>The Rent Management System Team</p>
        
        <div class="footer">
            <p><a href="https://rent.ragodevs.com">rent.ragodevs.com</a></p>
            <p>© 2025 Rent Management System. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`

// enableCORS middleware configures CORS settings for Echo
func enableCORS() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	})
}

// ipRestriction middleware checks if the client's IP is in the allowed list
func ipRestriction(allowedIPs string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			clientIP := c.RealIP()

			// If no allowed IPs are specified, allow all requests
			if allowedIPs == "" {
				return next(c)
			}

			// Split the allowed IPs string into a slice
			allowedIPList := strings.Split(allowedIPs, ",")

			// Check if the client IP is in the allowed list
			allowed := false
			for _, ip := range allowedIPList {
				if strings.TrimSpace(ip) == clientIP {
					allowed = true
					break
				}
			}

			// If the client IP is not allowed, return 403 Forbidden
			if !allowed {
				slog.Info(fmt.Sprintf("Blocked request from unauthorized IP: %s", clientIP))
				return c.JSON(http.StatusForbidden, map[string]string{"error": "Access denied"})
			}

			return next(c)
		}
	}
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

func (mc *MailConfig) sendWelcomeEmail(data SignupData) error {

	tmpl, err := template.New("email").Parse(welcome_template)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %v", err)
	}

	var emailBody bytes.Buffer
	if err := tmpl.Execute(&emailBody, data); err != nil {
		return fmt.Errorf("failed to execute email template: %v", err)
	}

	recipients := []string{data.Email}

	subject := "Welcome to Rent Management System - Account Activation Required"
	msg := fmt.Sprintf("Subject: %s\nTo: %s\nContent-Type: text/html\n\n%s", subject, data.Email, emailBody.String())

	auth := smtp.PlainAuth("", mc.USER, mc.PWD, mc.HOST)
	err = smtp.SendMail(mc.HOST+":"+mc.PORT, auth, mc.USER, recipients, []byte(msg))
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

func (mc *MailConfig) sendActivateEmail(data ActivateOrResetData) error {

	tmpl, err := template.New("email").Parse(activate_template)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %v", err)
	}

	var emailBody bytes.Buffer
	if err := tmpl.Execute(&emailBody, data); err != nil {
		return fmt.Errorf("failed to execute email template: %v", err)
	}

	recipients := []string{data.Email}

	subject := "Rent Management System - Account Activation Required"
	msg := fmt.Sprintf("Subject: %s\nTo: %s\nContent-Type: text/html\n\n%s", subject, data.Email, emailBody.String())

	auth := smtp.PlainAuth("", mc.USER, mc.PWD, mc.HOST)
	err = smtp.SendMail(mc.HOST+":"+mc.PORT, auth, mc.USER, recipients, []byte(msg))
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

func (mc *MailConfig) sendPasswordResetEmail(data ActivateOrResetData) error {

	tmpl, err := template.New("email").Parse(pwdreset_template)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %v", err)
	}

	var emailBody bytes.Buffer
	if err := tmpl.Execute(&emailBody, data); err != nil {
		return fmt.Errorf("failed to execute email template: %v", err)
	}

	recipients := []string{data.Email}

	subject := "Password Reset Request for Rent Management System"
	msg := fmt.Sprintf("Subject: %s\nTo: %s\nContent-Type: text/html\n\n%s", subject, data.Email, emailBody.String())

	auth := smtp.PlainAuth("", mc.USER, mc.PWD, mc.HOST)
	err = smtp.SendMail(mc.HOST+":"+mc.PORT, auth, mc.USER, recipients, []byte(msg))
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

func (mc *MailConfig) sendResetCompletedEmail(data ResetCompleteData) error {

	tmpl, err := template.New("email").Parse(completedreset_template)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %v", err)
	}

	var emailBody bytes.Buffer
	if err := tmpl.Execute(&emailBody, data); err != nil {
		return fmt.Errorf("failed to execute email template: %v", err)
	}

	recipients := []string{data.Email}

	subject := "Password Changed for Rent Management System"
	msg := fmt.Sprintf("Subject: %s\nTo: %s\nContent-Type: text/html\n\n%s", subject, data.Email, emailBody.String())

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
	flag.StringVar(&mc.RECEPIENTS, "RECEPIENTS", os.Getenv("RECEPIENTS"), "RECEPIENTS")
	flag.StringVar(&mc.ALLOWED_IP, "ALLOWED_IP", os.Getenv("ALLOWED_IP"), "ALLOWED_IP")

	flag.Parse()

	// Create a new Echo instance
	e := echo.New()

	// Set up validator
	e.Validator = &CustomValidator{validator: validator.New()}

	// Set up middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(enableCORS())
	e.Use(ipRestriction(mc.ALLOWED_IP))

	// Handle contact form submission
	e.POST("/submit-contact", func(c echo.Context) error {
		var form ContactForm
		if err := validateRequest(c, &form); err != nil {
			return err
		}

		recipients := strings.Split(mc.RECEPIENTS, ",")

		if err := mc.sendEmail(form, recipients); err != nil {
			log.Printf("Error sending email: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to send emails"})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "Emails sent successfully!"})
	})

	// Handle signup
	e.POST("/rent-signup", func(c echo.Context) error {
		var data SignupData
		if err := validateRequest(c, &data); err != nil {
			return err
		}

		if err := mc.sendWelcomeEmail(data); err != nil {
			log.Printf("Error sending email: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to send email"})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "Email sent successfully!"})
	})

	// Handle account activation
	e.POST("/rent-activate", func(c echo.Context) error {
		var data ActivateOrResetData
		if err := validateRequest(c, &data); err != nil {
			return err
		}

		if err := mc.sendActivateEmail(data); err != nil {
			log.Printf("Error sending email: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to send email"})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "Email sent successfully!"})
	})

	// Handle password reset request
	e.POST("/rent-resetpwd", func(c echo.Context) error {
		var data ActivateOrResetData
		if err := validateRequest(c, &data); err != nil {
			return err
		}

		if err := mc.sendPasswordResetEmail(data); err != nil {
			log.Printf("Error sending email: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to send email"})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "Email sent successfully!"})
	})

	// Handle completed password reset
	e.POST("/rent-completedpwdreset", func(c echo.Context) error {
		var data ResetCompleteData
		if err := validateRequest(c, &data); err != nil {
			return err
		}

		if err := mc.sendResetCompletedEmail(data); err != nil {
			log.Printf("Error sending email: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to send email"})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "Email sent successfully!"})
	})

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	slog.Info(fmt.Sprintf("Server is running on port %s", port))
	e.Logger.Fatal(e.Start(":" + port))
}
