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

type ArticleComment struct {
	Id        uint64    `json:"id" gorm:"primary_key"`
	ParentId  uint64    `json:"parent_id" gorm:"index;default:null;"`
	UserId    uint64    `json:"user_id" gorm:"index;not null"`
	ArticleId uint64    `json:"article_id" gorm:"index;not null"`
	Comment   string    `json:"comment"  gorm:"type:longtext;not null"`
	Status    uint8     `json:"status" gorm:"index;default:0"`
	CreatedAt time.Time `gorm:"index;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"index;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (ArticleComment) TableName() string {
	return "articles_comments"
}
