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

package seeds

import (
	"backend/models"
	"backend/utilities"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"github.com/Pallinder/go-randomdata"
	"github.com/bxcodec/faker/v4"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/kristijorgji/goseeder"
	math "math/rand"
	"strconv"
	"time"
)

func init() {
	goseeder.Register(CreateUser)
	goseeder.Register(CreateReference)
	goseeder.Register(CreateContact)
	goseeder.Register(CreateCustomer)
	goseeder.Register(CreateFaq)
	goseeder.Register(CreateService)
	goseeder.Register(CreateSlider)
	goseeder.Register(CreateTeam)
	goseeder.Register(CreateTestimonial)
	goseeder.Register(CreatePortfolio)
	goseeder.Register(CreateArticle)
}

func CreateUser(s goseeder.Seeder) {
	db := utilities.SetupDB()
	var totalRow int64
	db.Model(&models.User{}).Where("id <> 0").Count(&totalRow)

	bytes := make([]byte, 32) //generate a random 32 byte key for AES-256
	if _, err := rand.Read(bytes); err != nil {
		panic(err.Error())
	}

	key := hex.EncodeToString(bytes) //encode key in bytes to string and keep as secret, put in a vault
	encrypted := utilities.Encrypt("p4ssw0rd!", key)

	if totalRow == 0 {
		for i := 1; i <= 10; i++ {

			min := 1
			max := 2
			gender := math.Intn(max-min+1) + min
			firstName := ""
			genderChar := ""

			if gender == 1 {
				genderChar = "M"
				firstName = faker.FirstNameMale()
			} else {
				genderChar = "F"
				firstName = faker.FirstNameFemale()
			}

			user := models.User{
				FirstName:    firstName,
				LastName:     faker.LastName(),
				Gender:       genderChar,
				Country:      randomdata.Country(randomdata.FullCountry),
				Address:      sql.NullString{String: randomdata.Address(), Valid: true},
				AboutMe:      sql.NullString{String: randomdata.Paragraph(), Valid: true},
				Email:        randomdata.Email(),
				Phone:        randomdata.PhoneNumber(),
				Salt:		  key,
				Password:     encrypted,
				Status:       1,
				ConfirmToken: (uuid.New()).String(),
			}
			db.Create(&user)
		}
	}
}

func CreateReference(s goseeder.Seeder) {
	db := utilities.SetupDB()
	var totalRow int64
	db.Model(&models.ReferenceContent{}).Where("id <> 0").Count(&totalRow)
	if totalRow == 0 {

		var articles = [...]string{
			"Health and wellness",
			"Technology and gadgets",
			"Business and finance",
			"Travel and tourism",
			"Lifestyle and fashion",
		}

		var tags = [...]string{
			"Mental Health",
			"Fitness and Exercise",
			"Alternative Medicine",
			"Artificial Intelligence",
			"Network Security",
			"Cloud Computing",
			"Entrepreneurship",
			"Personal Finance",
			"Marketing and Branding",
			"Travel Tips and Tricks",
			"Cultural Experiences",
			"Destination Guides",
			"Beauty and Fashion Trends",
			"Celebrity News and Gossip",
			"Parenting and Family Life",
		}

		var portfolios = [...]string{
			"3D Modeling",
			"Web Application",
			"Mobile Application",
			"Illustrator Design",
			"UX Design",
		}

		for _, article := range articles {
			ar := models.ReferenceContent{
				Name:        article,
				Slug:        slug.Make(article),
				Description: randomdata.Paragraph(),
				Type:        1,
				Status:      1,
			}
			db.Create(&ar)
		}

		for _, tag := range tags {
			at := models.ReferenceContent{
				Name:        tag,
				Slug:        slug.Make(tag),
				Description: randomdata.Paragraph(),
				Type:        2,
				Status:      1,
			}
			db.Create(&at)
		}

		for _, portfolio := range portfolios {
			ap := models.ReferenceContent{
				Name:        portfolio,
				Slug:        slug.Make(portfolio),
				Description: randomdata.Paragraph(),
				Type:        3,
				Status:      1,
			}
			db.Create(&ap)
		}

	}
}

func CreateContact(s goseeder.Seeder) {
	db := utilities.SetupDB()
	var totalRow int64
	db.Model(&models.Contact{}).Where("id <> 0").Count(&totalRow)
	if totalRow == 0 {
		for i := 1; i <= 10; i++ {
			contact := models.Contact{
				Name:    faker.Name(),
				Email:   randomdata.Email(),
				Subject: faker.Sentence(),
				Message: faker.Paragraph(),
				Status:  0,
			}
			db.Create(&contact)
		}
	}
}

func CreateCustomer(s goseeder.Seeder) {
	db := utilities.SetupDB()
	var totalRow int64
	db.Model(&models.Customer{}).Where("id <> 0").Count(&totalRow)
	if totalRow == 0 {
		for i := 1; i <= 10; i++ {
			customer := models.Customer{
				Image:   ("customer" + strconv.Itoa(i) + ".jpg"),
				Name:    randomdata.City(),
				Email:   faker.Email(),
				Phone:   randomdata.PhoneNumber(),
				Address: randomdata.Address(),
				Sort:    uint16(i),
				Status:  1,
			}
			db.Create(&customer)
		}
	}
}

func CreateFaq(s goseeder.Seeder) {
	db := utilities.SetupDB()
	var totalRow int64
	db.Model(&models.Faq{}).Where("id <> 0").Count(&totalRow)
	if totalRow == 0 {
		for i := 1; i <= 10; i++ {
			faq := models.Faq{
				Question: faker.Sentence(),
				Answer:   randomdata.Paragraph(),
				Sort:     uint16(i),
				Status:   1,
			}
			db.Create(&faq)
		}
	}
}

func CreateService(s goseeder.Seeder) {
	db := utilities.SetupDB()
	var totalRow int64
	db.Model(&models.Service{}).Where("id <> 0").Count(&totalRow)
	if totalRow == 0 {
		var icons = [...]string{
			"bi bi-bicycle",
            "bi bi-bookmarks",
            "bi bi-box",
            "bi bi-building-add",
            "bi bi-calendar2-check",
            "bi bi-cart4",
            "bi bi-clipboard-data",
            "bi bi-gift",
            "bi bi-person-bounding-box",
		}
		for i, icon := range icons {
			service := models.Service{
				Icon:        icon,
				Title:       faker.Sentence(),
				Description: randomdata.Paragraph(),
				Status:      1,
				Sort:        uint16(i + 1),
			}
			db.Create(&service)
		}
	}
}

func CreateSlider(s goseeder.Seeder) {
	db := utilities.SetupDB()
	var totalRow int64
	db.Model(&models.Slider{}).Where("id <> 0").Count(&totalRow)
	if totalRow == 0 {
		for i := 1; i <= 5; i++ {
			slider := models.Slider{
				Image:       ("slider" + strconv.Itoa(i) + ".jpg"),
				Title:       faker.Sentence(),
				Description: randomdata.Paragraph(),
				Sort:        uint16(i),
				Status:      1,
			}
			db.Create(&slider)
		}
	}
}

func CreateTeam(s goseeder.Seeder) {
	db := utilities.SetupDB()
	var totalRow int64
	db.Model(&models.Team{}).Where("id <> 0").Count(&totalRow)
	if totalRow == 0 {
		for i := 1; i <= 10; i++ {
			team := models.Team{
				Image:        ("team" + strconv.Itoa(i) + ".jpg"),
				Name:         faker.Name(),
				Email:        randomdata.Email(),
				Phone:        randomdata.PhoneNumber(),
				PositionName: randomdata.SillyName(),
				Address:      randomdata.Address(),
				Twitter:      faker.Username(),
				Facebook:     faker.Username(),
				Instagram:    faker.Username(),
				LinkedIn:     faker.Username(),
				Sort:         uint16(i),
				Status:       1,
			}
			db.Create(&team)
		}
	}
}

func CreateTestimonial(s goseeder.Seeder) {
	db := utilities.SetupDB()
	var totalRow int64
	db.Model(&models.Testimonial{}).Where("id <> 0").Count(&totalRow)
	if totalRow == 0 {
		var customers []models.Customer
		db.Where("id <> 0").Find(&customers)
		for index, customer := range customers {
			testimonial := models.Testimonial{
				CustomerId: customer.Id,
				Image:      ("testimonial" + strconv.Itoa(index+1) + ".jpg"),
				Name:       faker.Name(),
				Position:   randomdata.SillyName(),
				Quote:      faker.Paragraph(),
				Sort:       uint16(index + 1),
				Status:     1,
			}
			db.Create(&testimonial)
		}
	}
}

func CreatePortfolio(s goseeder.Seeder) {
	db := utilities.SetupDB()
	var totalRow int64
	db.Model(&models.Portfolio{}).Where("id <> 0").Count(&totalRow)
	if totalRow == 0 {
		for i := 1; i <= 9; i++ {
			var category models.ReferenceContent
			var customer models.Customer
			db.Order("RAND()").Where("type = ?", 3).First(&category)
			db.Order("RAND()").Where("id <> 0").First(&customer)
			myDate, err := time.Parse("2006-01-02 15:04", faker.Date())
			portfolio := models.Portfolio{
				CustomerId:  customer.Id,
				ReferenceId: category.Id,
				Title:       faker.Sentence(),
				Description: randomdata.Paragraph(),
				ProjectUrl:  faker.URL(),
				Sort:        uint16(i),
				Status:      1,
			}

			if err == nil {
				portfolio.ProjectDate = sql.NullTime{Time: myDate}
			}

			db.Create(&portfolio)

			for j := 1; j <= 4; j++ {
				status := 0
				if j == 1 {
					status = 1
				}
				portfolioImage := models.PortfolioImage{
					PortfolioId: portfolio.Id,
					Image:       ("portfolio" + strconv.Itoa(j+1) + ".jpg"),
					Status:      uint8(status),
				}
				db.Create(&portfolioImage)
			}

		}
	}
}

func CreateArticle(s goseeder.Seeder) {
	db := utilities.SetupDB()
	var totalRow int64
	db.Model(&models.Article{}).Where("id <> 0").Count(&totalRow)
	if totalRow == 0 {
		var users []models.User
		db.Where("id <> 0").Find(&users)
		for index, user := range users {

			title := faker.Sentence()
			slug := slug.Make(title)
			article := models.Article{
				UserId:      user.Id,
				Image:       ("article" + strconv.Itoa(index+1) + ".jpg"),
				Title:       title,
				Slug:        slug,
				Description: faker.Sentence(),
				Content:     randomdata.Paragraph(),
				Status:      uint8(1),
			}
			db.Create(&article)

			var categories []models.ReferenceContent
			var tags []models.ReferenceContent
			var comments []models.User

			db.Order("RAND()").Where("type = 1").Limit(3).Find(&categories)
			db.Order("RAND()").Where("type = 2").Limit(5).Find(&tags)
			db.Order("RAND()").Where("id != " + strconv.Itoa(int(user.Id))).Limit(2).Find(&comments)

			for _, category := range categories {
				articleCategories := models.ArticleReference{
					ArticleId:   article.Id,
					ReferenceId: category.Id,
				}
				db.Create(&articleCategories)
			}

			for _, tag := range tags {
				articleTags := models.ArticleReference{
					ArticleId:   article.Id,
					ReferenceId: tag.Id,
				}
				db.Create(&articleTags)
			}

			for _, comment := range comments {
				articleComment := models.ArticleComment{
					ArticleId: article.Id,
					UserId:    comment.Id,
					Comment:   randomdata.Paragraph(),
					Status:    uint8(1),
				}
				db.Create(&articleComment)
			}

		}
	}
}
