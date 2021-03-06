package models

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dimitrijed93/glossary/api/util"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Entry struct {
	Id            uint   `json:"id" gorm:"primaryKey"`
	Definition    string `json:"definition"`
	FolderID      int    `json:"folderId"`
	SourceID      int
	DestinationID uint
	Source        EntryItem `json:"source" gorm:"association_foreignkey:Id"`
	Destination   EntryItem `json:"destination" gorm:"association_foreignkeyId"`
}

func NewEntry() *Entry {
	return &Entry{}
}

func (entry *Entry) Url() string {
	return "/entries"
}

func (entry *Entry) DetailUrl() string {
	return "/entries/id"
}

func (entry *Entry) GetAll(db *gorm.DB, v *util.Validator) func(c echo.Context) error {

	return func(c echo.Context) error {
		var entries []Entry
		db.Find(&entries)

		return c.JSON(http.StatusOK, entries)
	}
}

func (entry *Entry) GetByFolder(db *gorm.DB, v *util.Validator) func(c echo.Context) error {

	return func(c echo.Context) error {
		var result []EntryDto
		param := c.Param("id")
		folderId, err := strconv.Atoi(param)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "id is required")
		}

		db.Table("entries").Select("definition", "Id").Where("folder_id = ?", folderId).Scan(&result)

		return c.JSON(http.StatusOK, result)
	}
}

func (entry *Entry) GetById(db *gorm.DB, v *util.Validator) func(c echo.Context) error {

	return func(c echo.Context) error {
		var entry Entry
		entryId := c.Param("id")
		db.Model(&entry).Where("id=?", entryId)

		return c.JSON(http.StatusOK, entry)
	}

}

func (entry *Entry) Create(db *gorm.DB, v *util.Validator) func(c echo.Context) error {

	return func(c echo.Context) error {
		en := NewEntry()

		if err := c.Bind(en); err != nil {
			return fmt.Errorf("id is mandatory for update")
		}

		if err := v.Validate(en); err != nil {
			return err
		}

		db.Save(en)
		return c.JSON(http.StatusCreated, en)
	}

}

func (entry *Entry) Update(db *gorm.DB, v *util.Validator) func(c echo.Context) error {

	return func(c echo.Context) error {
		en := new(Type)

		if err := c.Bind(en); err != nil {
			return err
		}

		itemId := c.Param("id")

		if itemId == "0" {
			return v.InvalidId()
		}

		if err := v.Validate(en); err != nil {
			return err
		}

		db.Model(&en).Where("id=?", itemId).Updates(&en)

		return c.JSON(http.StatusOK, en)
	}

}

func (entry *Entry) Delete(db *gorm.DB, v *util.Validator) func(c echo.Context) error {

	return func(c echo.Context) error {
		itemId := c.Param("id")

		en := &Entry{}

		db.First(&en, itemId)

		if itemId == "0" {
			return v.InvalidId()
		}

		if en != nil {
			db.Delete(&en)
		}

		return c.JSON(http.StatusNoContent, en)

	}
}
