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

package utilities

import (
	"backend/models"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"os"
)

func SetupDB() *gorm.DB {
	godotenv.Load(".env")
	USER := os.Getenv("DB_USERNAME")
	PASS := os.Getenv("DB_PASSWORD")
	HOST := os.Getenv("DB_HOST")
	PORT := os.Getenv("DB_PORT")
	DBNAME := os.Getenv("DB_DATABASE")
	URL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", USER, PASS, HOST, PORT, DBNAME)
	db, err := gorm.Open(os.Getenv("DB_CONNECTION"), URL)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func Config() {
	db := SetupDB()
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(models.User{})
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(models.Contact{})
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(models.Customer{})
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(models.Faq{})
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(models.ReferenceContent{})
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(models.Service{})
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(models.Slider{})
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(models.Team{})
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(models.Article{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(models.ArticleComment{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(models.ArticleComment{}).AddForeignKey("article_id", "articles(id)", "RESTRICT", "RESTRICT")
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(models.ArticleComment{}).AddForeignKey("parent_id", "articles_comments(id)", "RESTRICT", "RESTRICT")
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(models.ArticleReference{}).AddForeignKey("article_id", "articles(id)", "RESTRICT", "RESTRICT")
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(models.ArticleReference{}).AddForeignKey("reference_id", "references_contents(id)", "RESTRICT", "RESTRICT")
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(models.Portfolio{}).AddForeignKey("reference_id", "references_contents(id)", "RESTRICT", "RESTRICT")
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(models.PortfolioImage{}).AddForeignKey("portfolio_id", "portfolios(id)", "RESTRICT", "RESTRICT")
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(models.Testimonial{}).AddForeignKey("customer_id", "customers(id)", "RESTRICT", "RESTRICT")
}
