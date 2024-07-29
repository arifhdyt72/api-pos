package model

import (
	"fmt"
	"os"
	"strings"
	"test_backend_esb/tools"

	"gorm.io/gorm"
)

type Category struct {
	MavisModel
	Name        string `json:"name"`
	StoreID     uint   `json:"store_id"`
	Store       *Store `json:"store"`
	Icon        string `json:"icon"`
	SourceImage string `json:"source_image" gorm:"-"`
}

func (i *Category) BeforeSave(tx *gorm.DB) (err error) {
	if i.SourceImage != "" {
		if i.Icon != "" {
			err := os.Remove(i.Icon)
			if err != nil {
				return err
			}
		}
		filename := strings.ReplaceAll(i.Name, " ", "-")
		path, err := tools.SaveImageToDisk(i.SourceImage, filename)
		if err != nil {
			return err
		}

		i.Icon = fmt.Sprintf("/%s", path)
		i.SourceImage = ""
	}

	return nil
}

func (i *Category) BeforeDelete(tx *gorm.DB) (err error) {
	if err := tx.First(i, i.ID).Error; err != nil {
		return err
	}

	if i.Icon != "" {
		err := os.Remove(i.Icon)
		if err != nil {
			return err
		}
	}

	return nil
}
