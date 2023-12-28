package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/go-mail/mail"
	_ "github.com/go-sql-driver/mysql"
	"github.com/vanng822/go-premailer/premailer"

	"github.com/TulioH/correos-masivos/src/config"
)

// Define a command-line flag
var createTable = flag.Bool("create-table", false, "Set to true to create the table")

var insertFromExcel = flag.Bool("insert-excel", false, "Set to true to insert from excel")
var sendFromExcel = flag.Bool("send-from-excel", false, "Set to true to send email from excel")
var excelpath string

var timeout = 100 * time.Second

func main() {
	// Set up flag to accept command-line argument for excel path
	flag.StringVar(&excelpath, "p", "default", "a string variable")
	// Parse command-line flags
	flag.Parse()

	// Load environment variables
	env := config.NewEnv()

	// Set up email variables
	from := env.Email
	password := env.Password
	smtpHost := env.SMTPHost
	smtpPort := env.SMTPPort
	subject := env.Subject

	count := env.Begin

	// Set up logger flags
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Create SQL connection string
	sqlConnectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", env.DBUser, env.DBPass, env.DBHost, env.DBPort, env.DBName)

	// Open database connection
	db, err := sql.Open("mysql", sqlConnectionString)
	if err != nil {
		log.Println(err)
		time.Sleep(timeout * time.Second)
		os.Exit(1)
	}
	defer db.Close()

	// Check if the create-table flag is set to true, create the table and exit
	if *createTable {
		migrate(db)
		time.Sleep(timeout * time.Second)
		os.Exit(1)
	}

	// Check if the insert-from-excel flag is set to true, insert rows from excel and exit
	if *insertFromExcel {
		if err := InsertRowsFromExcel(db, excelpath); err != nil {
			log.Println(err)
			time.Sleep(timeout * time.Second)
			os.Exit(1)
		}
		time.Sleep(timeout * time.Second)
		os.Exit(1)
	}

	if *sendFromExcel {
		SendEmailFromExcel(excelpath, env, from, password, smtpHost, smtpPort, subject)
		time.Sleep(timeout * time.Second)
		os.Exit(1)
	}

	// Construct query to fetch emails from the database
	query := fmt.Sprintf(`SELECT NOMBRE, CORREO FROM correo limit %d, %d`, env.Begin, env.EmailsForPack)

	// Execute the query
	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		time.Sleep(timeout * time.Second)
		os.Exit(1)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			NOMBRE string
			CORREO string
		)

		// Scan the values from the database into variables
		if err := rows.Scan(&NOMBRE, &CORREO); err != nil {
			log.Println(err)
			time.Sleep(timeout * time.Second)
			os.Exit(1)
		}

		// Print the name and email for debugging
		fmt.Printf("NOMBRE: %s, CORREO: %s\n", NOMBRE, CORREO)

		// Send email
		SendEmail(env, from, password, []string{CORREO}, smtpHost, smtpPort, subject, NOMBRE)

		// Increment count
		count++
		fmt.Println(count)
	}

	// Check for any errors during iteration
	if err = rows.Err(); err != nil {
		log.Println(err)
		time.Sleep(timeout * time.Second)
		os.Exit(1)
	}

	time.Sleep(timeout * time.Second)
}

// SendEmail sends an email using the provided SMTP settings and email content.
//
// Parameters:
// - env: the environment configuration.
// - from: the sender's email address.
// - password: the sender's email password.
// - to: a list of recipient email addresses.
// - smtpHost: the SMTP server host.
// - smtpPort: the SMTP server port.
// - subject: the email subject.
// - NOMBRE: the name to use in the email body.
func SendEmail(env *config.Env, from string, password string, to []string, smtpHost string, smtpPort int, subject string, NOMBRE string) {

	htmlContent, err := os.ReadFile(env.EmailBody)
	if err != nil {
		log.Println(err)
		time.Sleep(timeout * time.Second)
		os.Exit(1)
	}
	modifiedHtmlContent := strings.Replace(string(htmlContent), "{NOMBRE_PROVEEDOR}", NOMBRE, -1)
	prem, err := premailer.NewPremailerFromString(modifiedHtmlContent, nil)
	if err != nil {
		log.Println(err)
		time.Sleep(timeout * time.Second)
		os.Exit(1)
	}
	body, err := prem.Transform()
	if err != nil {
		log.Println(err)
		time.Sleep(timeout * time.Second)
		os.Exit(1)
	}

	m := mail.NewMessage()

	m.SetHeader("From", from)

	m.SetHeader("To", to...)

	// m.SetAddressHeader("Cc", "oliver.doe@example.com", "Oliver")

	m.SetHeader("Subject", subject)

	m.SetBody("text/html", body)

	if _, err := os.Stat("./email/logo.png"); err == nil {
		m.Embed("./email/logo.png", mail.SetHeader(map[string][]string{"Content-ID": {fmt.Sprintf("<%s>", "logo")}}))
	}

	attachmentsPath := env.Attachments

	files, err := os.ReadDir(attachmentsPath)
	if err != nil {
		log.Println(err)
	}

	for _, nameFile := range files {
		m.Attach(attachmentsPath + "/" + nameFile.Name())
	}

	d := mail.NewDialer(smtpHost, smtpPort, from, password)

	// Send the email to Kate, Noah and Oliver.

	if err := d.DialAndSend(m); err != nil {
		log.Println(err)
		time.Sleep(timeout * time.Second)
		os.Exit(1)
	}

	fmt.Println("Email sent!")

}

// migrate creates a table named "correo" if it doesn't exist in the database specified by the db parameter.
//
// db: a pointer to a sql.DB object representing the database connection.
func migrate(db *sql.DB) {

	// SQL statement to create a table
	createTableSQL := `CREATE TABLE IF NOT EXISTS correo (
		id INT AUTO_INCREMENT,
		NOMBRE VARCHAR(255) NOT NULL,
		CORREO VARCHAR(255) NOT NULL,
		PRIMARY KEY (id)
	);`

	tx, err := db.Begin()
	if err != nil {
		log.Printf("error starting transaction: %v", err)
		time.Sleep(timeout * time.Second)
		os.Exit(1)
	}

	stmt, err := tx.Prepare(createTableSQL)
	if err != nil {
		log.Printf("error preparing statement: %v", err)
		time.Sleep(timeout * time.Second)
		os.Exit(1)
	}
	defer stmt.Close()

	// Execute the SQL statement to create the table
	_, err = stmt.Exec()
	if err != nil {
		log.Println("Error creating table:", err)
		time.Sleep(timeout * time.Second)
		os.Exit(1)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		log.Printf("error committing transaction: %v", err)
		time.Sleep(timeout * time.Second)
		os.Exit(1)
	}

	log.Println("Table created successfully")
	time.Sleep(timeout * time.Second)
	os.Exit(1)
}

// insertRowsFromExcel opens an Excel file at the given filePath,
// reads all the rows from the first sheet, and inserts them into
// a SQL database using the provided database handle db.
// This function assumes that the first row of the sheet contains headers
// and therefore starts inserting from the second row.
// It uses a transaction to ensure that all inserts are successful,
// and commits the transaction at the end. If an error occurs during
// any step, it returns an error and the transaction is rolled back.
func InsertRowsFromExcel(db *sql.DB, filePath string) error {
	// Open the Excel file
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("error opening Excel file: %v", err)
	}

	// Get all rows from the first sheet
	rows := f.GetRows(f.GetSheetName(1))

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}

	// Prepare the SQL statement for inserting rows into the table
	stmt, err := tx.Prepare("INSERT INTO correo (NOMBRE, CORREO) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmt.Close()

	// Iterate over the rows, skipping the first one if it contains headers
	for _, row := range rows[1:] {
		// Execute the insert statement for each row
		_, err := stmt.Exec(row[0], row[1])
		if err != nil {
			return fmt.Errorf("error executing insert statement for row: %v, error: %v", row, err)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return fmt.Errorf("error committing transaction: %v", err)
	}

	// Log a success message
	log.Println("All rows have been inserted successfully.")
	return nil
}

func SendEmailFromExcel(filePath string, env *config.Env, from string, Password string, smtpHost string, smtpPort int, subject string) {
	// Open the Excel file
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		log.Printf("error opening Excel file: %v", err)
		time.Sleep(timeout * time.Second)
		os.Exit(1)
	}

	// Get all rows from the first sheet
	rows := f.GetRows(f.GetSheetName(1))
	count := env.Begin
	fmt.Println("cantidad de filas", len(rows))
	// Iterate over the rows, skipping the first one if it contains headers
	for _, row := range rows[(count + 1):] {
		fmt.Printf("NOMBRE: %s, CORREO: %s\n", row[0], row[1])
		SendEmail(env, from, Password, []string{row[1]}, smtpHost, smtpPort, subject, row[0])
		count++
		fmt.Println(count)
	}

	// Log a success message
	log.Println("All rows have been sended successfully.")
}

func LogFatal(err error) {
	log.Println(err)
	time.Sleep(timeout * time.Second)
	os.Exit(1)
}
