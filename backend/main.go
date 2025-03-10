package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

	"github.com/jung-kurt/gofpdf"
)

type Student struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Course struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type SelectedCourse struct {
	StudentID int `json:"student_id"`
	CourseID  int `json:"course_id"`
}

type Certificate struct {
	StudentID int    `json:"student_id"`
	CourseID  int    `json:"course_id"`
	PDFPath   string `json:"pdf_path"`
}

var db *sql.DB
var jwtSecret []byte

func connectDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Database unreachable", err)
	}
	fmt.Println("Connected to database")

	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
}

func signup(c *gin.Context) {
	var student Student
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(student.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	_, err = db.Exec("INSERT INTO students (email, password) VALUES (?, ?)", student.Email, string(hashedPassword))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Signup successful"})
}

func login(c *gin.Context) {
	var student Student
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var storedPassword string
	err := db.QueryRow("SELECT id, password FROM students WHERE email = ?", student.Email).Scan(&student.ID, &storedPassword)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(student.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": student.Email,
		"exp":   time.Now().Add(time.Hour * 1).Unix(),
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	// dynamically generate token and send it to frontend as response along with student ID
	// ✅ Now returning student ID in response
	c.JSON(http.StatusOK, gin.H{
		"token":     tokenString,
		"studentId": student.ID, // Send student ID to frontend
	})
}

// ---------------all courses list-- API Handlers -----------------//
func getCourses(c *gin.Context) {
	rows, err := db.Query("SELECT id, name FROM courses")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var courses []Course
	for rows.Next() {
		var course Course
		if err := rows.Scan(&course.ID, &course.Name); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning data"})
			return
		}
		courses = append(courses, course)
	}
	c.JSON(http.StatusOK, courses)
}

// ---------------select course-- API Handlers -----------------//
func SelectCourse(c *gin.Context) {
	var selection struct {
		StudentID int `json:"student_id"`
		CourseID  int `json:"course_id"`
	}
	if err := c.ShouldBindJSON(&selection); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Check if course is already selected
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM selected_courses WHERE student_id = ? AND course_id = ?", selection.StudentID, selection.CourseID).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if count > 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Course already selected!"})
		return
	}

	// Insert course selection
	_, err = db.Exec("INSERT INTO selected_courses (student_id, course_id) VALUES (?, ?)", selection.StudentID, selection.CourseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to select course"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course Selected! Choose Another Course."})
}

// api to get selected courses by student id
func GetSelectedCourses(c *gin.Context) {
	studentID := c.Param("student_id")

	rows, err := db.Query("SELECT courses.id, courses.name FROM selected_courses JOIN courses ON selected_courses.course_id = courses.id WHERE selected_courses.student_id = ?", studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var selectedCourses []map[string]interface{}
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning data"})
			return
		}
		selectedCourses = append(selectedCourses, gin.H{
			"id":   id,
			"name": name,
		})
	}

	c.JSON(http.StatusOK, selectedCourses)
}

func main() {
	connectDB()
	router := gin.Default()

	// CORS Middleware Configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // React frontend का origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Routes
	router.POST("/signup", signup)
	router.POST("/login", login)
	router.GET("/courses", getCourses)
	router.POST("/select-course", SelectCourse)
	router.GET("/selected-courses/:student_id", GetSelectedCourses)
	// router.GET("/certificate/:student_id/:course_id", generateCertificate)
	router.Static("/uploads", "./uploads")
	router.GET("/generate-certificate/:student_id/:course_id", generateCertificate)

	// Start background cleanup job
	go cleanupCertificates()
	// Start the server
	log.Fatal(router.Run(":8080"))
}

func generateCertificate(c *gin.Context) {
	studentID := c.Param("student_id")
	courseID := c.Param("course_id")

	// Fetch student name and course name
	query := `SELECT students.email, courses.name 
			FROM students 
			JOIN selected_courses ON students.id = selected_courses.student_id 
			JOIN courses ON courses.id = selected_courses.course_id 
			WHERE students.id = ? AND courses.id = ?`

	var studentEmail, courseName string
	err := db.QueryRow(query, studentID, courseID).Scan(&studentEmail, &courseName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student or Course not found"})
		return
	}

	// Check if certificate already exists
	var existingURL string
	err = db.QueryRow("SELECT certificate_url FROM certificates WHERE student_id = ? AND course_id = ?", studentID, courseID).Scan(&existingURL)
	if err == nil {
		// Return the existing certificate
		c.JSON(http.StatusOK, gin.H{
			"message":      "Certificate already exists",
			"student_id":   studentID,
			"student_name": studentEmail, // Returning email as the student identifier
			"course_name":  courseName,
			"pdf_url":      fmt.Sprintf("http://localhost:8080/%s", existingURL),
		})
		return
	}

	// Ensure uploads directory exists
	uploadDir := "uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.Mkdir(uploadDir, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create uploads directory"})
			return
		}
	}

	// Define PDF file path
	fileName := fmt.Sprintf("certificate_%s_%s.pdf", studentID, courseID)
	filePath := filepath.Join(uploadDir, fileName)

	// Generate PDF certificate
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Certificate of Completion")
	pdf.Ln(20)
	pdf.SetFont("Arial", "", 12)
	pdf.MultiCell(0, 10, fmt.Sprintf("This certifies that %s has successfully completed the course '%s'.", studentEmail, courseName), "0", "C", false)

	// Save the PDF file
	if err := pdf.OutputFileAndClose(filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PDF"})
		return
	}

	// Store the file path in the database
	_, err = db.Exec("INSERT INTO certificates (student_id, course_id, certificate_url) VALUES (?, ?, ?)", studentID, courseID, filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save certificate info"})
		return
	}

	// Respond with success
	c.JSON(http.StatusOK, gin.H{
		"message":      "Certificate generated successfully",
		"student_id":   studentID,
		"student_name": studentEmail, // Using email instead of name
		"course_name":  courseName,
		"pdf_url":      fmt.Sprintf("http://localhost:8080/%s", filePath), // Full URL for easy frontend access
	})
}

func cleanupCertificates() {
	for {
		rows, err := db.Query("SELECT id, certificate_url FROM certificates")
		if err != nil {
			log.Println("Error fetching certificates:", err)
			continue
		}

		var deletedIDs []int

		for rows.Next() {
			var certID int
			var certPath string
			err := rows.Scan(&certID, &certPath)
			if err != nil {
				log.Println("Error scanning row:", err)
				continue
			}

			// Check if file exists
			if _, err := os.Stat(certPath); os.IsNotExist(err) {
				// File is missing, delete from database
				_, err := db.Exec("DELETE FROM certificates WHERE id = ?", certID)
				if err != nil {
					log.Println("Error deleting certificate record:", err)
				} else {
					log.Println("Deleted missing certificate record:", certID)
					deletedIDs = append(deletedIDs, certID) // Store deleted IDs
				}
			}
		}

		rows.Close()

		// ✅ Rearrange IDs sequentially if any row was deleted
		if len(deletedIDs) > 0 {
			_, err := db.Exec(`
				SET @num := 0;
				UPDATE certificates SET id = @num := @num + 1;
				ALTER TABLE certificates AUTO_INCREMENT = 1;
			`)
			if err != nil {
				log.Println("Error rearranging certificate IDs:", err)
			} else {
				log.Println("Certificate IDs rearranged sequentially")
			}
		}

		// Run cleanup every 10 seconds
		time.Sleep(10 * time.Second)
	}
}
