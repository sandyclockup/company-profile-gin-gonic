/**
 * This file is part of the Sandy Andryanto Company Profile Website.
 *
 * @author     Sandy Andryanto <sandy.andryanto.dev@gmail.com>
 * @copyright  2024
 *
 * For the full copyright and license information,
 * please view the LICENSE.md file that was distributed
 * with this source code.
 */

package controllers

import (
	"backend/models"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/bxcodec/faker/v4"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type ArticleResult struct {
	Id          uint64    `json:"id"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Image       string    `json:"image"`
	Gender      string    `json:"gender"`
	AboutMe     string    `json:"about_me"`
	Categories  string    `json:"categories"`
}

type TestimonialResult struct {
	Id           uint64    `json:"id"`
	CustomerName string    `json:"customer_name"`
	Image        string    `json:"image"`
	Name         string    `json:"name"`
	Position     string    `json:"position"`
	Quote        string    `json:"quote"`
	Sort         uint16    `json:"sort"`
	Status       uint8     `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type FormMessage struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type FormSubscribe struct {
	Email string `json:"email"`
}

func PagePing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "ok", "status": true, "data": nil})
}

func PageHome(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)

	var articles []ArticleResult
	db.Raw(`
		SELECT 
			a.id,
			a.title,
			a.slug,
			a.description,
			a.created_at,
			u.first_name,
			u.last_name,
			u.image,
			u.gender,
			u.about_me,
			(SELECT 
				GROUP_CONCAT(r.name SEPARATOR ',') AS r 
				FROM references_contents r
				WHERE r.id IN (
					SELECT reference_id
					FROM articles_references
					WHERE article_id = a.id
				)
			) as categories
		FROM articles a
		INNER JOIN users u ON u.id = a.user_id
		WHERE a.status = 1 ORDER BY RAND() LIMIT 3 
	`).Scan(&articles)

	var testimonials []TestimonialResult
	db.Raw(`SELECT t.*, c.name customer_name FROM testimonials t INNER JOIN customers c ON c.id = t.customer_id WHERE t.status = 1 ORDER BY RAND() LIMIT 1`).Scan(&testimonials)

	header := map[string]string{
		"title":       faker.Sentence(),
		"description": randomdata.Paragraph(),
	}

	var sliders []models.Slider
	db.Where("status = 1").Order("sort").Find(&sliders)

	var services []models.Service
	db.Where("status = 1").Limit(4).Order("RAND()").Find(&services)

	c.JSON(http.StatusOK, gin.H{"header": header, "sliders": sliders, "services": services, "testimonial": testimonials[0], "articles": articles})
}

func PageAbout(c *gin.Context) {

	header := map[string]string{
		"title":       faker.Sentence(),
		"description": randomdata.Paragraph(),
	}

	section1 := map[string]string{
		"title":       faker.Sentence(),
		"description": randomdata.Paragraph(),
	}

	section2 := map[string]string{
		"title":       faker.Sentence(),
		"description": randomdata.Paragraph(),
	}

	var teams []models.Team
	db := c.MustGet("db").(*gorm.DB)
	db.Where("status = 1").Order("RAND()").Find(&teams)

	c.JSON(http.StatusOK, gin.H{"header": header, "section1": section1, "section2": section2, "teams": teams})
}

func PageService(c *gin.Context) {

	header := map[string]string{
		"title":       faker.Sentence(),
		"description": randomdata.Paragraph(),
	}

	db := c.MustGet("db").(*gorm.DB)

	var services []models.Service
	var customers []models.Customer
	var testimonials []models.Testimonial

	db.Where("status = 1").Order("RAND()").Find(&services)
	db.Where("status = 1").Order("RAND()").Find(&customers)
	db.Where("status = 1").Order("RAND()").Find(&testimonials)

	c.JSON(http.StatusOK, gin.H{"header": header, "services": services, "customers": customers, "testimonials": testimonials})
}

func PageFaq(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)
	var faq1 []models.Faq
	var faq2 []models.Faq

	db.Where("status = 1 AND sort <= 5").Order("RAND()").Find(&faq1)
	db.Where("status = 1 AND sort > 5").Order("RAND()").Find(&faq2)

	c.JSON(http.StatusOK, gin.H{"faq1": faq1, "faq2": faq2})
}

func PageContact(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)

	var services []models.Service
	db.Where("status = 1").Limit(4).Order("RAND()").Find(&services)

	c.JSON(http.StatusOK, gin.H{"message": "ok", "status": true, "services": services})
}

func PageMessage(c *gin.Context) {

	var input FormMessage
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(strings.TrimSpace(input.Email)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The email field is required.!"})
		return
	}

	if len(strings.TrimSpace(input.Name)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The name field is required.!"})
		return
	}

	if len(strings.TrimSpace(input.Subject)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The subject field is required.!"})
		return
	}

	if len(strings.TrimSpace(input.Message)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The message field is required.!"})
		return
	}

	db := c.MustGet("db").(*gorm.DB)

	contact := models.Contact{
		Email:   input.Email,
		Name:    input.Name,
		Subject: input.Subject,
		Message: input.Message,
	}
	db.Create(&contact)

	c.JSON(http.StatusOK, gin.H{"message": "ok", "status": true, "data": contact})
}

func PageSubscribe(c *gin.Context) {

	var input FormSubscribe
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(strings.TrimSpace(input.Email)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The email field is required.!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok", "status": true, "data": nil})
}

func PageGetFile(c *gin.Context) {
	path := os.Getenv("UPLOAD_PATH")
	filepath := path + "/" + c.Query("param")
	c.File(filepath)
}
