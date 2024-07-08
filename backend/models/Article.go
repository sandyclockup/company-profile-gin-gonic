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

package models

import "time"

type Article struct {
	Id          uint64    `json:"id" gorm:"primary_key"`
	UserId      uint64    `json:"user_id" gorm:"index;not null"`
	Image       string    `json:"image" gorm:"index;size:191;default:null;"`
	Title       string    `json:"title" gorm:"index;size:255;not null"`
	Slug        string    `json:"slug" gorm:"index;size:255;not null"`
	Description string    `json:"description" gorm:"index;size:255;default:null;"`
	Content     string    `json:"content"  gorm:"type:longtext;default:null;"`
	Status      uint8     `json:"status" gorm:"index;default:0"`
	CreatedAt   time.Time `gorm:"index;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"index;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (Article) TableName() string {
	return "articles"
}
