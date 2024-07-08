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
	"net/http"
	"strings"
	"strconv"
	"time"
	"backend/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/dgrijalva/jwt-go"
)

type ArticleResultDetail struct {
	Id          uint64    `json:"id"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	Content 	string    `json:"content"`
	CreatedAt   time.Time `json:"created_at"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Image       string    `json:"image"`
	Gender      string    `json:"gender"`
	AboutMe     string    `json:"about_me"`
	Categories  string    `json:"categories"`
}
type ArticleTreeModel struct {
	Id  		uint64	  			`json:"id"`
	ParentId  	uint64    			`json:"parent_id"`
	Comment 	string    			`json:"comment"`
	CreatedAt   time.Time 			`json:"created_at"`
	FirstName   string    			`json:"first_name"`
	LastName    string    			`json:"last_name"`
	Gender      string    			`json:"gender"`
}

type ArticleTreeResultModel struct {
	Id  		uint64	  			
	ParentId  	uint64    			
	Comment 	string    			
	CreatedAt   time.Time 			
	FirstName   string    			
	LastName    string    			
	Gender      string    			
	Childern	[]ArticleTreeResultModel
}

type FormComment struct {
	Comment string `json:"comment"`
}

func ArticleList(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)

	page := 1
	if len(strings.TrimSpace(c.Query("page"))) > 0 {
		page_, err := strconv.ParseInt(c.Query("page"), 0, 32)
		if err == nil {
			page = int(page_)
		}
	}

	limit := 3 * page

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
		WHERE a.status = 1 ORDER BY id DESC
	`).Scan(&articles)

	var news_articles []ArticleResult
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
		WHERE a.status = 1 ORDER BY a.id DESC LIMIT 3 OFFSET 1
	`).Scan(&news_articles)

	var stories []ArticleResult
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
		WHERE a.status = 1 ORDER BY a.id DESC LIMIT ?
	`, limit).Scan(&stories)
	
	var total_article int
	db.Model(&models.Article{}).Where("status = 1").Count(&total_article)

	c.JSON(http.StatusOK, gin.H{ "continue": limit <= total_article, "new_article": articles[0], "news_articles": news_articles, "page": page, "stories": stories })
}

func ArticleDetail(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)

	var articles []ArticleResultDetail
	db.Raw(`
		SELECT 
			a.id,
			a.title,
			a.slug,
			a.description,
			a.content,
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
		WHERE a.slug = ? LIMIT 1
	`, c.Param("slug")).Scan(&articles)

	if len(articles) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok", "status": true, "data": articles[0]})
}

func ArticleCommentCreate(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)
	var input FormComment
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(strings.TrimSpace(input.Comment)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The comment field is required.!"})
		return
	}

	auth := c.MustGet("claims").(jwt.MapClaims)
	var user models.User
	if err := db.Where("id = ?", auth["id"]).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found"})
		return
	}

	article_id, _ := strconv.Atoi(c.Param("id"))

	comment := models.ArticleComment{
		UserId:		user.Id,
		ArticleId:	uint64(article_id),
		Comment: 	input.Comment,
		Status: 	1,
	}
	db.Create(&comment)

	c.JSON(http.StatusOK, gin.H{"message": "ok", "status": true, "data": comment})
}

func ArticleCommentList(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var result []ArticleTreeModel
	db.Raw(`
		SELECT 
			c.id,
			IFNULL(c.parent_id, 0) parent_id,
			c.comment,
			c.created_at,
			u.first_name,
			u.last_name,
			u.gender
		FROM articles_comments c
		INNER JOIN users u ON u.id = c.user_id
		WHERE c.article_id = ?
		ORDER BY c.id DESC
	`, c.Param("id")).Scan(&result)
	commentTree := ArticleTree(result, 0)
	c.JSON(http.StatusOK, gin.H{"message": "ok", "status": true, "data": commentTree})
}

func ArticleTree(elements []ArticleTreeModel, ParentId uint64) []ArticleTreeResultModel{
	var branch []ArticleTreeResultModel
	for _, element := range elements {

		if(element.ParentId == ParentId){

			childern := []ArticleTreeResultModel{}

			getChildren := ArticleTree(elements, element.Id)
			if(len(getChildren) > 0){
				childern = getChildren
			}

			var obj = ArticleTreeResultModel{
				Id: element.Id,
				ParentId: element.ParentId,
				Comment: element.Comment,
				CreatedAt: element.CreatedAt,
				FirstName: element.FirstName,
				LastName: element.LastName,
				Gender: element.Gender,
				Childern: childern,
			}

			branch = append(branch, obj)
		}
	}
	return branch
}
