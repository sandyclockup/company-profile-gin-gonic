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
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type PortfolioResult struct {
	Id           uint64       `json:"id"`
	CustomerId   uint64       `json:"customer_id" `
	CustomerName string       `json:"customer_name"`
	ReferenceId  uint64       `json:"reference_id"`
	CategoryName string       `json:"category_name"`
	Title        string       `json:"title"`
	Description  string       `json:"description"`
	ProjectDate  sql.NullTime `json:"project_date"`
	ProjectUrl   string       `json:"project_url"`
	Sort         uint16       `json:"sort"`
	Status       uint8        `json:"status"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

func PortfolioList(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var portfolios []PortfolioResult
	db.Raw(`
		SELECT p.*, 
			c.name customer_name, 
			r.name category_name 
		FROM portfolios p
		INNER JOIN customers c ON c.id = p.customer_id 
		INNER JOIN references_contents r ON r.id = p.reference_id
		WHERE p.status = 1
		ORDER BY p.id DESC
	`).Scan(&portfolios)
	c.JSON(http.StatusOK, gin.H{"message": "ok", "status": true, "data": portfolios})
}

func PortfolioDetail(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)

	var portfolios []PortfolioResult

	db.Raw(`
		SELECT p.*, 
			c.name customer_name, 
			r.name category_name 
		FROM portfolios p
		INNER JOIN customers c ON c.id = p.customer_id 
		INNER JOIN references_contents r ON r.id = p.reference_id
		WHERE p.id = ` + c.Param("id") + `
	`).Scan(&portfolios)

	var images []models.PortfolioImage
	db.Where("portfolio_id = " + c.Param("id")).Order("id").Find(&images)

	c.JSON(http.StatusOK, gin.H{"portfolio": portfolios[0], "images": images})
}
