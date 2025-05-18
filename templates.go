package main

// Common CSS styles for all email templates
const commonStyles = `
    body {
        font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
        line-height: 1.6;
        color: #333;
        max-width: 600px;
        margin: 0 auto;
        background-color: #f9f9f9;
    }
    .container {
        background-color: #ffffff;
        border-radius: 8px;
        box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
        overflow: hidden;
        margin: 20px auto;
    }
    .header {
        background-color: #4361ee;
        padding: 25px 20px;
        text-align: center;
        color: white;
    }
    .header h2 {
        margin: 0;
        font-weight: 600;
        font-size: 24px;
    }
    .header p {
        margin: 10px 0 0;
        opacity: 0.9;
    }
    .content {
        padding: 30px 25px;
    }
    .field {
        margin-bottom: 20px;
        padding-bottom: 15px;
        border-bottom: 1px solid #eee;
    }
    .field:last-child {
        border-bottom: none;
    }
    .label {
        font-weight: 600;
        color: #4361ee;
        display: inline-block;
        min-width: 120px;
    }
    .message-box {
        background-color: #f8f9fa;
        border-left: 4px solid #4361ee;
        padding: 15px;
        margin-top: 15px;
        border-radius: 0 4px 4px 0;
    }
    .button {
        background-color: #4361ee;
        color: white !important;    
        padding: 12px 28px;
        text-decoration: none;
        border-radius: 4px;
        font-weight: 600;
        display: inline-block;
        margin: 20px 0;
        text-align: center;
        transition: all 0.3s ease;
    }
    .button:hover {
        background-color: #3a56d4;
        transform: translateY(-2px);
        box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
    }
    .important-note {
        background-color: #fff8e1;
        border-left: 4px solid #ffc107;
        padding: 15px;
        margin: 20px 0;
        font-size: 0.95em;
        border-radius: 0 4px 4px 0;
    }
    .footer {
        margin-top: 30px;
        padding: 20px;
        background-color: #f8f9fa;
        font-size: 0.9em;
        color: #666;
        text-align: center;
        border-top: 1px solid #eee;
        border-radius: 0 0 8px 8px;
    }
    .contact-info {
        margin-top: 25px;
        padding: 15px;
        background-color: #f1f3f9;
        border-radius: 6px;
    }
    a {
        color: #4361ee;
        text-decoration: none;
    }
    a:hover {
        text-decoration: underline;
    }
    .logo {
        max-width: 150px;
        margin-bottom: 10px;
    }
`

const emailTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        ${commonStyles}
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h2>New Contact Form Submission</h2>
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
                <span class="label">Service:</span> {{.Service}}
            </div>
            
            <div class="field">
                <span class="label">Message:</span>
                <div class="message-box">{{.Message}}</div>
            </div>
            
            <div class="footer">
                <p>This is an automated message from your website contact form.</p>
            </div>
        </div>
    </div>
</body>
</html>
`

const welcome_template = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        ${commonStyles}
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h2>Welcome to Rent Management System</h2>
        </div>
        
        <div class="content">
            <p>Thank you for choosing Rent Management System for your property management needs. We're delighted to welcome you to our platform.</p>
            
            <div class="field">
                <p><strong>Your account has been created successfully with the following details:</strong></p>
                <p><strong class="label">User ID:</strong> {{.ID}}</p>
            </div>
            
            <h3 style="color: #4361ee; margin-top: 30px;">Important: Please Activate Your Account</h3>
            
            <p>To complete your registration and access all features of our platform, please activate your account by clicking the button below:</p>
            
            <div style="text-align: center;">
                <a href="https://rent.ragodevs.com/activate?token={{.Token}}" class="button">Activate Account</a>
            </div>
            
            <div class="important-note">
                <p>Please note that this activation link will expire in 3 days and can only be used once.</p>
            </div>
            
            <div class="contact-info">
                <p><strong>Need help?</strong> Contact our support team:</p>
                <p>Email: <a href="mailto:support@ragodevs.com">support@ragodevs.com</a><br>
                Phone: +255 654 051 622</p>
            </div>
            
            <p>We look forward to helping you streamline your property management operations.</p>
            
            <p>Best regards,<br>
            The Rent Management System Team</p>
        </div>
        
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
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        ${commonStyles}
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h2>Account Activation Required</h2>
        </div>
        
        <div class="content">
            <h3 style="color: #4361ee; margin-top: 0;">Important: Please Activate Your Account</h3>
            
            <p>Your account requires activation to continue using our services. Please click the button below to activate:</p>
            
            <div style="text-align: center;">
                <a href="https://rent.ragodevs.com/activate?token={{.Token}}" class="button">Activate Account</a>
            </div>
            
            <div class="important-note">
                <p><strong>Security Note:</strong> This activation link will expire in 3 days and can only be used once.</p>
            </div>
            
            <div class="contact-info">
                <p><strong>Need help?</strong> Contact our support team:</p>
                <p>Email: <a href="mailto:support@ragodevs.com">support@ragodevs.com</a><br>
                Phone: +255 654 051 622</p>
            </div>
            
            <p>We look forward to helping you streamline your property management operations.</p>
            
            <p>Best regards,<br>
            The Rent Management System Team</p>
        </div>
        
        <div class="footer">
            <p><a href="https://rent.ragodevs.com">rent.ragodevs.com</a></p>
            <p>© 2025 Rent Management System. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`

const pwdreset_template = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        ${commonStyles}
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h2>Password Reset Request</h2>
        </div>
        
        <div class="content">
            <p>We received a request to reset your password for your Rent Management System account. If you did not make this request, you can safely ignore this email.</p>
            
            <p>To reset your password, please click the button below:</p>
            
            <div style="text-align: center;">
                <a href="https://rent.ragodevs.com/reset?token={{.Token}}" class="button">Reset Password</a>
            </div>
            
            <div class="important-note">
                <p><strong>Security Note:</strong> This reset link will expire in 45 minutes and can only be used once.</p>
            </div>
            
            <div class="contact-info">
                <p><strong>Need help?</strong> Contact our support team:</p>
                <p>Email: <a href="mailto:support@ragodevs.com">support@ragodevs.com</a><br>
                Phone: +255 654 051 622</p>
            </div>
            
            <p>Best regards,<br>
            The Rent Management System Team</p>
        </div>
        
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
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        ${commonStyles}
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h2>Password Changed</h2>
        </div>
        
        <div class="content">
            <div style="text-align: center; margin-bottom: 20px;">
                <span style="font-size: 48px; color: #4CAF50;">✓</span>
                <h3 style="color: #4CAF50; margin-top: 10px;">Success!</h3>
                <p>You have successfully changed your password for Rent Management System.</p>
            </div>
            
            <div class="important-note">
                <p><strong>Security Alert:</strong> If this wasn't done by you, please immediately reset your password and contact our support team.</p>
            </div>
            
            <div class="contact-info">
                <p><strong>Need help?</strong> Contact our support team:</p>
                <p>Email: <a href="mailto:support@ragodevs.com">support@ragodevs.com</a><br>
                Phone: +255 654 051 622</p>
            </div>
            
            <p>Best regards,<br>
            The Rent Management System Team</p>
        </div>
        
        <div class="footer">
            <p><a href="https://rent.ragodevs.com">rent.ragodevs.com</a></p>
            <p>© 2025 Rent Management System. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`
