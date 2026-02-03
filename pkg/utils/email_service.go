package utils

import (
	"fmt"
	"net/smtp"

	"go.uber.org/zap"
)

// EmailService menangani pengiriman email
type EmailService struct {
	logger *zap.Logger
	host   string
	port   string
	sender string
	pass   string
}

// NewEmailService membuat instance baru EmailService
func NewEmailService(logger *zap.Logger, smtpConfig SMTPConfig) *EmailService {
	return &EmailService{
		logger: logger,
		host:   smtpConfig.Host,
		port:   smtpConfig.Port,
		sender: smtpConfig.Email,
		pass:   smtpConfig.Password,
	}
}

// SendOTP mengirim OTP ke email
func (es *EmailService) SendOTP(toEmail, otpCode, purpose string) error {
	if es.host == "" || es.port == "" || es.sender == "" || es.pass == "" {
		es.logger.Warn("SMTP not configured, skipping email send", zap.String("to_email", toEmail))
		return nil // Skip jika SMTP tidak dikonfigurasi
	}

	subject := "Your OTP Code"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #007bff; color: white; padding: 20px; text-align: center; border-radius: 5px 5px 0 0; }
        .content { background-color: #f8f9fa; padding: 20px; border-radius: 0 0 5px 5px; }
        .otp-code { font-size: 32px; font-weight: bold; color: #007bff; text-align: center; margin: 20px 0; letter-spacing: 5px; }
        .expiry { color: #666; font-size: 14px; text-align: center; margin-top: 20px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>OTP Verification</h1>
        </div>
        <div class="content">
            <p>Hello,</p>
            <p>Your OTP code for %s is:</p>
            <div class="otp-code">%s</div>
            <p class="expiry">This code will expire in 10 minutes. Do not share this code with anyone.</p>
            <p>If you didn't request this, please ignore this email.</p>
            <p>Best regards,<br>POS Application Team</p>
        </div>
    </div>
</body>
</html>
    `, purpose, otpCode)

	// Setup SMTP
	addr := fmt.Sprintf("%s:%s", es.host, es.port)
	auth := smtp.PlainAuth("", es.sender, es.pass, es.host)

	// Headers
	headers := fmt.Sprintf("MIME-version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\nFrom: %s\r\nTo: %s\r\nSubject: %s\r\n",
		es.sender, toEmail, subject)

	message := []byte(headers + "\r\n" + body)

	// Send email
	err := smtp.SendMail(addr, auth, es.sender, []string{toEmail}, message)
	if err != nil {
		es.logger.Error("Failed to send OTP email",
			zap.String("to_email", toEmail),
			zap.Error(err),
		)
		return fmt.Errorf("failed to send OTP email: %w", err)
	}

	es.logger.Info("OTP email sent successfully",
		zap.String("to_email", toEmail),
		zap.String("purpose", purpose),
	)

	return nil
}

// SendPasswordResetEmail mengirim link reset password ke email
func (es *EmailService) SendPasswordResetEmail(toEmail, resetToken string) error {
	if es.host == "" || es.port == "" || es.sender == "" || es.pass == "" {
		es.logger.Warn("SMTP not configured, skipping email send", zap.String("to_email", toEmail))
		return nil
	}

	subject := "Password Reset Request"
	resetLink := fmt.Sprintf("http://your-app-domain.com/reset-password?token=%s", resetToken)
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #007bff; color: white; padding: 20px; text-align: center; border-radius: 5px 5px 0 0; }
        .content { background-color: #f8f9fa; padding: 20px; border-radius: 0 0 5px 5px; }
        .button { display: inline-block; background-color: #007bff; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px; margin: 20px 0; }
        .expiry { color: #666; font-size: 14px; margin-top: 20px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Password Reset Request</h1>
        </div>
        <div class="content">
            <p>Hello,</p>
            <p>You have requested to reset your password. Click the button below to proceed:</p>
            <a href="%s" class="button">Reset Password</a>
            <p class="expiry">This link will expire in 1 hour. If you didn't request this, please ignore this email.</p>
            <p>Best regards,<br>POS Application Team</p>
        </div>
    </div>
</body>
</html>
    `, resetLink)

	addr := fmt.Sprintf("%s:%s", es.host, es.port)
	auth := smtp.PlainAuth("", es.sender, es.pass, es.host)

	headers := fmt.Sprintf("MIME-version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\nFrom: %s\r\nTo: %s\r\nSubject: %s\r\n",
		es.sender, toEmail, subject)

	message := []byte(headers + "\r\n" + body)

	err := smtp.SendMail(addr, auth, es.sender, []string{toEmail}, message)
	if err != nil {
		es.logger.Error("Failed to send password reset email",
			zap.String("to_email", toEmail),
			zap.Error(err),
		)
		return fmt.Errorf("failed to send password reset email: %w", err)
	}

	es.logger.Info("Password reset email sent successfully",
		zap.String("to_email", toEmail),
	)

	return nil
}

// SendWelcomeEmail mengirim email selamat datang ke user baru
func (es *EmailService) SendWelcomeEmail(toEmail, name, tempPassword string) error {
	if es.host == "" || es.port == "" || es.sender == "" || es.pass == "" {
		es.logger.Warn("SMTP not configured, skipping email send", zap.String("to_email", toEmail))
		return nil
	}

	subject := "Welcome to POS Application"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #28a745; color: white; padding: 20px; text-align: center; border-radius: 5px 5px 0 0; }
        .content { background-color: #f8f9fa; padding: 20px; border-radius: 0 0 5px 5px; }
        .credentials { background-color: #e9ecef; padding: 15px; border-radius: 5px; margin: 20px 0; }
        .credentials p { margin: 10px 0; }
        .warning { color: #dc3545; font-weight: bold; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Welcome to POS Application</h1>
        </div>
        <div class="content">
            <p>Hello %s,</p>
            <p>Your account has been successfully created. Here are your login credentials:</p>
            <div class="credentials">
                <p><strong>Email:</strong> %s</p>
                <p><strong>Temporary Password:</strong> %s</p>
                <p class="warning">⚠️ Please change your password after first login.</p>
            </div>
            <p>If you have any questions, please contact our support team.</p>
            <p>Best regards,<br>POS Application Team</p>
        </div>
    </div>
</body>
</html>
    `, name, toEmail, tempPassword)

	addr := fmt.Sprintf("%s:%s", es.host, es.port)
	auth := smtp.PlainAuth("", es.sender, es.pass, es.host)

	headers := fmt.Sprintf("MIME-version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\nFrom: %s\r\nTo: %s\r\nSubject: %s\r\n",
		es.sender, toEmail, subject)

	message := []byte(headers + "\r\n" + body)

	err := smtp.SendMail(addr, auth, es.sender, []string{toEmail}, message)
	if err != nil {
		es.logger.Error("Failed to send welcome email",
			zap.String("to_email", toEmail),
			zap.Error(err),
		)
		return fmt.Errorf("failed to send welcome email: %w", err)
	}

	es.logger.Info("Welcome email sent successfully",
		zap.String("to_email", toEmail),
	)

	return nil
}
