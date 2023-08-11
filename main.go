package main

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/labstack/echo/v4"
)

func sendSMS(mobno, message string) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"), // Change this to your desired AWS region
	})
	if err != nil {
		return err
	}

	svc := sns.New(sess)
	params := &sns.PublishInput{
		Message:     aws.String(message),
		PhoneNumber: aws.String(mobno),
	}

	_, err = svc.Publish(params)
	return err
}

func sendEmail(mailID, message string) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"), // Change this to your desired AWS region
	})
	if err != nil {
		return err
	}

	svc := ses.New(sess)
	params := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(mailID),
			},
		},
		Message: &ses.Message{
			Subject: &ses.Content{
				Data: aws.String("Subject"),
			},
			Body: &ses.Body{
				Text: &ses.Content{
					Data: aws.String(message),
				},
			},
		},
		Source: aws.String("keerthanashanmugam252000@gmail.com"), // Change this to your verified SES email
	}

	_, err = svc.SendEmail(params)
	return err
}

func sendSMSAndEmail(c echo.Context) error {
	var req struct {
		Email   string `json:"email"`
		MobNo   string `json:"mobno"`
		Message string `json:"message"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request format")
	}

	// Send SMS
	if err := sendSMS(req.MobNo, req.Message); err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to send SMS: %v", err))
	}

	// Send Email
	if err := sendEmail(req.Email, req.Message); err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to send email: %v", err))
	}

	return c.JSON(http.StatusOK, "SMS and Email sent successfully")
}

func main() {
	e := echo.New()

	e.POST("/send-sms-email", sendSMSAndEmail)

	e.Start(":8080")
}
