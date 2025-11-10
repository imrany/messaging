# mailer Package

The `mailer` package provides functionality for sending emails,
including OTP (One-Time Password) emails, using SMTP.

## Configuration

The package uses the `SMTPConfig` struct to store SMTP server configuration details.

```go
type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Email    string
}
```

- `Host`: The SMTP server hostname.
- `Port`: The SMTP server port.
- `Username`: The SMTP username for authentication.
- `Password`: The SMTP password for authentication.
- `Email`: The email address used as the sender.

## Usage

### Sending a Basic Email

The `SendEmail` function sends a generic email.

```go
type EmailData struct {
	To      []string
	Subject string
	Body    string
	IsHTML  bool
}

func SendEmail(emailData EmailData, config SMTPConfig) error
```

- `emailData`: An `EmailData` struct containing the recipient(s), subject, body, and a flag indicating whether the body is HTML.
- `config`: An `SMTPConfig` struct containing the SMTP server configuration.

**Example:**

```go
package main

import (
	"fmt"
	"log"
	"net/smtp"
	"github.com/your-username/whats-email/pkg/mailer" // Replace with the actual import path
)

func main() {
	config := mailer.SMTPConfig{
		Host:     "smtp.example.com",
		Port:     587,
		Username: "your_username",
		Password: "your_password",
		Email:    "your_email@example.com",
	}

	emailData := mailer.EmailData{
		To:      []string{"recipient@example.com"},
		Subject: "Test Email",
		Body:    "This is a test email.",
		IsHTML:  false,
	}

	err := mailer.SendEmail(emailData, config)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Email sent successfully!")
}

```

### Sending an OTP Email

The `SendOTP` function generates and sends an OTP email.

```go
type OtpPurpose string

const (
	OtpPurposeLogin         OtpPurpose = "login"
	OtpPurposePasswordReset OtpPurpose = "password_reset"
	OtpPurposeVerification  OtpPurpose = "verification"
	OtpPurposeRegistration  OtpPurpose = "registration"
)


func SendOTP(email string, purpose OtpPurpose, otp string, config SMTPConfig) (string, error)
```

- `email`: The recipient's email address.
- `purpose`: The purpose of the OTP (e.g., "login", "password_reset").
- `otp`: The OTP code.
- `config`: An `SMTPConfig` struct containing the SMTP server configuration.

**Example:**

```go
package main

import (
	"fmt"
	"log"
	"github.com/your-username/whats-email/pkg/mailer" // Replace with the actual import path
)

func main() {
	config := mailer.SMTPConfig{
		Host:     "smtp.example.com",
		Port:     587,
		Username: "your_username",
		Password: "your_password",
		Email:    "your_email@example.com",
	}

	email := "recipient@example.com"
	purpose := mailer.OtpPurposeLogin
	otp := "123456"

	sentOTP, err := mailer.SendOTP(email, purpose, otp, config)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("OTP %s sent successfully to %s!\n", sentOTP, email)
}
```

### Sending an OTP Email with a Custom Template

The `SendOTPWithCustomTemplate` function allows you to use custom HTML and text templates for the OTP email.

```go
func SendOTPWithCustomTemplate(email string, purpose OtpPurpose, otp string, subject, htmlTemplate, textTemplate string, config SMTPConfig) (string, error)
```

- `email`: The recipient's email address.
- `purpose`: The purpose of the OTP.
- `otp`: The OTP code.
- `subject`: The email subject.
- `htmlTemplate`: The HTML template for the email body. Use `{{OTP}}` and `{{PURPOSE}}` placeholders.
- `textTemplate`: The text template for the email body. Use `{{OTP}}` and `{{PURPOSE}}` placeholders.
- `config`: An `SMTPConfig` struct containing the SMTP server configuration.

### Resending an OTP Email

The `SendOTPWithExistingCode` function resends an OTP using an existing OTP code.

```go
func SendOTPWithExistingCode(email string, purpose OtpPurpose, otp string, config SMTPConfig) error
```

## OTP Data Handling

The `OTPData` struct represents OTP information:

```go
type OTPData struct {
	Code      string
	ExpiresAt time.Time
	Purpose   OtpPurpose
}
```

- `Code`: The OTP code.
- `ExpiresAt`: The expiration time of the OTP.
- `Purpose`: The purpose of the OTP.

### OTP Expiration Duration

The `GetOTPExpirationDuration` function returns the default expiration duration for each OTP purpose.

### Creating OTP Data

The `CreateOTPData` function creates an `OTPData` struct with the appropriate expiration time.

### Checking OTP Expiration

The `IsOTPExpired` method checks if an OTP has expired.

### Checking OTP Validity for a Specific Purpose

The `IsValidForPurpose` method checks if an OTP is valid for a specific purpose.
