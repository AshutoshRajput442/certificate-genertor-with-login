package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
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

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	router.POST("/signup", signup)
	router.POST("/login", login)
	router.GET("/courses", getCourses)
	router.POST("/select-course", SelectCourse)
	router.GET("/selected-courses/:student_id", GetSelectedCourses)
	// http.HandleFunc("/select-course", SelectCourse)
	// router.POST("/generate-certificate", generateCertificate)

	log.Fatal(router.Run(":8080"))
}
