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

type ArticleReference struct {
	ArticleId   uint64 `gorm:"primary_key;auto_increment:false"`
	ReferenceId uint64 `gorm:"primary_key;auto_increment:false"`
}

func (ArticleReference) TableName() string {
	return "articles_references"
}
