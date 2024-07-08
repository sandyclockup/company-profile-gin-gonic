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

type Team struct {
	Id           uint64    `json:"id" gorm:"primary_key"`
	Image        string    `json:"image" gorm:"index;size:191;default:null;"`
	Name         string    `json:"name" gorm:"index;size:255;not null"`
	Email        string    `json:"email" gorm:"index;size:255;not null"`
	Phone        string    `json:"phone" gorm:"index;size:255;not null"`
	PositionName string    `json:"position_name" gorm:"index;size:255;not null"`
	Address      string    `json:"address" gorm:"type:text;default:null;"`
	Twitter      string    `json:"twitter" gorm:"index;size:191;default:null;"`
	Facebook     string    `json:"facebook" gorm:"index;size:191;default:null;"`
	Instagram    string    `json:"instagram" gorm:"index;size:191;default:null;"`
	LinkedIn     string    `json:"linked_in" gorm:"index;size:191;default:null;"`
	Sort         uint16    `json:"sort" gorm:"index;default:0"`
	Status       uint8     `json:"status" gorm:"index;default:0"`
	CreatedAt    time.Time `gorm:"index;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time `gorm:"index;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (Team) TableName() string {
	return "teams"
}
