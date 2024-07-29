package model

import (
	"os"
	"strings"
	"test_backend_esb/tools"

	"gorm.io/gorm"
)

type Item struct {
	MavisModel
	Name        string    `json:"name"`
	CategoryID  uint      `json:"category_id"`
	Category    *Category `json:"category"`
	StoreID     uint      `json:"store_id"`
	Store       *Store    `json:"store"`
	Stock       int       `json:"stock"`
	BasePrice   int64     `json:"base_price"`
	Detail      string    `json:"detail"`
	Photo       string    `json:"photo"`
	SourceImage string    `json:"source_image" gorm:"-"`
}

func (i *Item) BeforeSave(tx *gorm.DB) (err error) {
	if i.SourceImage != "" {
		if i.Photo != "" {
			err := os.Remove(i.Photo)
			if err != nil {
				return err
			}
		}
		filename := strings.ReplaceAll(i.Name, " ", "-")
		path, err := tools.SaveImageToDisk(i.SourceImage, filename)
		if err != nil {
			return err
		}

		i.Photo = path
		i.SourceImage = ""
	}

	return nil
}

func (i *Item) BeforeDelete(tx *gorm.DB) (err error) {
	if err := tx.First(i, i.ID).Error; err != nil {
		return err
	}

	if i.Photo != "" {
		err := os.Remove(i.Photo)
		if err != nil {
			return err
		}
	}

	return nil
}
