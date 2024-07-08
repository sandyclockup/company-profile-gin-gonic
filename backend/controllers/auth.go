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
	service "backend/services"
	"backend/utilities"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterInput struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"password_confirm"`
}

type ForgotEmailInput struct {
	Email string `json:"email"`
}

type ResetEmailInput struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"password_confirm"`
}

func AuthLogin(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)

	var user models.User

	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(strings.TrimSpace(input.Email)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The email field is required.!"})
		return
	}

	if len(strings.TrimSpace(input.Password)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The password field is required.!"})
		return
	}

	if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user " + input.Email + " not found!"})
		return
	}

	decrypt := utilities.Decrypt(user.Password, user.Salt)

	if user.Status == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You need to confirm your account. We have sent you an activation code, please check your email.!"})
		return
	}

	if input.Password != decrypt {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect password!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": service.JWTAuthService().GenerateToken(int(user.Id), user.Email, true)})
}

func AuthRegister(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)

	var user models.User

	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(strings.TrimSpace(input.Email)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The email field is required.!"})
		return
	}

	if len(strings.TrimSpace(input.Password)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The passwword field is required.!"})
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

	if err := db.Where("email = ?", input.Email).First(&user).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The email already exists"})
		return
	}

	bytes := make([]byte, 32) //generate a random 32 byte key for AES-256
	if _, err := rand.Read(bytes); err != nil {
		panic(err.Error())
	}

	key := hex.EncodeToString(bytes) //encode key in bytes to string and keep as secret, put in a vault
	encrypted := utilities.Encrypt(input.Password, key)
	token := utilities.Encrypt(input.Email, key)

	User := models.User{
		Email:        input.Email,
		Password:     encrypted,
		Status:       1,
		ConfirmToken: token,
		Salt:         key,
	}

	db.Create(&User)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func AuthConfirm(c *gin.Context) {

	var user models.User

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("confirm_token = ? AND status = ? ", c.Param("token"), 0).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	db.Exec("UPDATE users SET status = 1, updated_at = NOW() WHERE confirm_token = ? ", c.Param("token"))
	c.JSON(http.StatusOK, gin.H{"data": "User has been confirmed"})
}

func AuthEmailForgot(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)

	var user models.User

	var input ForgotEmailInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(strings.TrimSpace(input.Email)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The email field is required.!"})
		return
	}

	if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "We can't find a user with that e-mail address."})
		return
	}

	bytes := make([]byte, 32) //generate a random 32 byte key for AES-256
	if _, err := rand.Read(bytes); err != nil {
		panic(err.Error())
	}

	key := hex.EncodeToString(bytes) //encode key in bytes to string and keep as secret, put in a vault
	token := utilities.Encrypt(input.Email, key)

	db.Exec("UPDATE users SET reset_token = ? WHERE email = ? ", token, input.Email)
	c.JSON(http.StatusOK, gin.H{"message": "We have e-mailed your password reset link!", "token": token})
}

func AuthEmailReset(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)

	var user models.User

	var input ResetEmailInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(strings.TrimSpace(input.Email)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The email field is required.!"})
		return
	}

	if len(strings.TrimSpace(input.Password)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The passwword field is required.!"})
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

	if err := db.Where("email = ? AND reset_token = ?", input.Email, c.Param("token")).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This password and email reset token is invalid."})
		return
	}

	bytes := make([]byte, 32) //generate a random 32 byte key for AES-256
	if _, err := rand.Read(bytes); err != nil {
		panic(err.Error())
	}

	key := hex.EncodeToString(bytes) //encode key in bytes to string and keep as secret, put in a vault
	encrypted := utilities.Encrypt(input.Password, key)

	var updatedInput models.User
	updatedInput.ResetToken = sql.NullString{String: "", Valid: true}
	updatedInput.Password = encrypted
	updatedInput.Salt = key
	db.Model(&user).Updates(updatedInput)

	c.JSON(http.StatusOK, gin.H{"message": "Your password has been reset!"})
}
