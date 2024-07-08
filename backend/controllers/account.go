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
	"backend/utilities"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

type FormProfile struct {
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Gender    string `json:"gender"`
	Country   string `json:"country"`
	Address   string `json:"address"`
	AboutMe   string `json:"about_me"`
}

type FormPassword struct {
	OldPassword     string `json:"old_password"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"password_confirm"`
}

func ProfileDetail(c *gin.Context) {

	auth := c.MustGet("claims").(jwt.MapClaims)
	db := c.MustGet("db").(*gorm.DB)

	var user models.User
	if err := db.Where("id = ?", auth["id"]).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok", "status": true, "data": user})
}

func ProfileUpdate(c *gin.Context) {

	authUser := c.MustGet("claims").(jwt.MapClaims)
	db := c.MustGet("db").(*gorm.DB)

	var user models.User

	var input FormProfile
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(strings.TrimSpace(input.Email)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The email field is required.!"})
		return
	}

	if err := db.Where("email = ? AND id != ?", input.Email, authUser["id"]).First(&user).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email address already exists !!"})
		return
	}

	if len(strings.TrimSpace(input.Phone)) > 0 {
		if err := db.Where("phone = ? AND id != ?", input.Phone, authUser["id"]).First(&user).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number already exists !!"})
			return
		}
	}

	if err := db.Where("id = ?", authUser["id"]).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	var updatedInput models.User
	updatedInput.Email = input.Email
	updatedInput.Phone = input.Phone
	updatedInput.FirstName = input.FirstName
	updatedInput.LastName = input.LastName
	updatedInput.Gender = input.Gender
	updatedInput.Country = input.Country
	updatedInput.Address = sql.NullString{String: input.Address, Valid: true}
	updatedInput.AboutMe = sql.NullString{String: input.AboutMe, Valid: true}
	db.Model(&user).Updates(updatedInput)

	c.JSON(http.StatusOK, gin.H{"message": "Your profile has been changed!", "status": true, "data": updatedInput})

}

func PasswordUpdate(c *gin.Context) {

	authUser := c.MustGet("claims").(jwt.MapClaims)
	db := c.MustGet("db").(*gorm.DB)

	var user models.User

	var input FormPassword
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(strings.TrimSpace(input.OldPassword)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The old_password field is required.!"})
		return
	}

	if len(strings.TrimSpace(input.Password)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The password field is required.!"})
		return
	}

	if len(strings.TrimSpace(input.ConfirmPassword)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The password_confirm field is required.!"})
		return
	}

	if len(strings.TrimSpace(input.Password)) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "you have to enter at least 8 digit!"})
		return
	}

	if strings.TrimSpace(input.Password) != strings.TrimSpace(input.ConfirmPassword) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "These passwords don't match!"})
		return
	}

	if err := db.Where("id = ?", authUser["id"]).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	decrypt := utilities.Decrypt(user.Password, user.Salt)

	if input.Password != decrypt {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect current password!"})
		return
	}

	bytes := make([]byte, 32) //generate a random 32 byte key for AES-256
	if _, err := rand.Read(bytes); err != nil {
		panic(err.Error())
	}

	key := hex.EncodeToString(bytes) //encode key in bytes to string and keep as secret, put in a vault
	encrypted := utilities.Encrypt(input.Password, key)

	var updatedInput models.User
	updatedInput.Password = encrypted
	updatedInput.Salt = key
	db.Model(&user).Updates(updatedInput)

	c.JSON(http.StatusOK, gin.H{"message": "Your password has been changed!"})

}

func ProfileUpload(c *gin.Context) {

	authUser := c.MustGet("claims").(jwt.MapClaims)
	db := c.MustGet("db").(*gorm.DB)
	var user models.User
	if err := db.Where("id = ?", authUser["id"]).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	godotenv.Load(".env")
	file, err := c.FormFile("file")
	currentTime := time.Now()
	datePath := currentTime.Format("2006-01-02")

	// The file cannot be received.
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "No file is received",
		})
		return
	}

	path := os.Getenv("UPLOAD_PATH")
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}

	realPath := path + "/" + datePath
	if _, err := os.Stat(realPath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(realPath, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}

	// Retrieve file information
	extension := filepath.Ext(file.Filename)
	// Generate random file name for the new uploaded file so it doesn't override the old file with same name
	newFileName := uuid.New().String() + extension

	// The file is received, so let's save it
	if err := c.SaveUploadedFile(file, realPath+"/"+newFileName); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to save the file",
		})
		return
	}

	result := datePath + "/" + newFileName

	if result != "" {

		if user.Image != "" {
			e := os.Remove(path + "/" + user.Image)
			if e != nil {
				log.Fatal(e)
			}
		}

		user.Image = result
		db.Model(&user).Updates(user)
	}

	// File saved successfully. Return proper result
	c.JSON(http.StatusOK, gin.H{"data": result})

}
