package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/smtp"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	from := "your-email@example.com"
	password := "your-password"
	to := []string{"recipient@example.com"}
	smtpHost := "smtp.example.com"
	smtpPort := "587"

	subject := "Subject: Your Subject Here\n"
	htmlContent, err := os.ReadFile("./correo.html")
	if err != nil {
		log.Fatal(err)
	}
	modifiedHtmlContent := strings.Replace(string(htmlContent), "{NOMBRE_PROVEEDOR}", "value_to_replace", -1)
	body := modifiedHtmlContent
	attachmentPath := "/path/to/attachment.pdf"

	// Create a new buffer to hold the message.
	var msg bytes.Buffer
	writer := multipart.NewWriter(&msg)

	// Set the headers for the email and the HTML part.
	headers := make(textproto.MIMEHeader)
	headers.Set("From", from)
	headers.Set("To", to[0])
	headers.Set("Subject", subject)
	headers.Set("MIME-Version", "1.0")
	headers.Set("Content-Type", "multipart/mixed; boundary="+writer.Boundary())
	writer.WriteField("Content-Type", "text/html; charset=\"utf-8\"")
	writer.WriteField("Content-Transfer-Encoding", "quoted-printable")

	for k, v := range headers {
		writer.WriteField(k, v[0])
	}

	// Write the HTML part of the email.
	htmlWriter, _ := writer.CreatePart(headers)
	htmlWriter.Write([]byte(body))

	// Attach the file.
	file, err := os.Open(attachmentPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	attachmentHeaders := textproto.MIMEHeader{}
	attachmentHeaders.Set("Content-Type", "application/octet-stream")
	attachmentHeaders.Set("Content-Transfer-Encoding", "base64")
	attachmentHeaders.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filepath.Base(attachmentPath)))

	attachmentWriter, _ := writer.CreatePart(attachmentHeaders)
	encodedFile := base64.StdEncoding.EncodeToString(fileBytes)
	attachmentWriter.Write([]byte(encodedFile))

	// Close the writer to finalize the message.
	writer.Close()

	// Set up authentication information.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Send the email.
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, msg.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Email sent successfully")
}
