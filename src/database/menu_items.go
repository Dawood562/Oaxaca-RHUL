package database

import (
	"errors"
	"fmt"
	"teamproject/database/models"
)

// AddItem adds the given item to the database.
// Returns an error if there is a problem adding the item.
// Item names must be unique.
func AddItem(item *models.MenuItem) error {
	result := db.Create(item)
	return result.Error
}

// EditItem edits the given item with new information.
// Returns an error if the item doesn't exist or the edit causes a conflict.
func EditItem(item *models.MenuItem) error {
	// Check that the item exists
	result := db.First(&models.MenuItem{ID: item.ID})
	if result.RowsAffected == 0 {
		return errors.New("item does not exist")
	}
	// Delete old allergen relations
	db.Where(&models.Allergen{ItemID: item.ID}).Delete(&models.Allergen{})
	// Update the item
	result = db.Save(&item)
	return result.Error
}

// RemoveItem removes an item from the menu with the given id
// Returns an error if the item could not be removed
func RemoveItem(id uint) error {
	result := db.Delete(&models.MenuItem{ID: id})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("item with id '%d' does not exist", id)
	}
	return nil
}

// QueryMenu returns the menu items from the database as a slice
// If filter is provided, the returned item slice will be filtered as such
func QueryMenu(filter *MenuFilter) []models.MenuItem {
	var data []models.MenuItem
	allergens := filterAllergens(filter.Allergens)
	dbLocal := db.Model(&models.MenuItem{}).Preload("Allergens")
	if len(filter.SearchTerm) > 3 {
		dbLocal = dbLocal.Where("LOWER(menu_items.name) LIKE LOWER(?)", fmt.Sprintf("%%%s%%", filter.SearchTerm))
	}
	if filter.MaxCalories > 0 {
		dbLocal = dbLocal.Where("calories <= ?", filter.MaxCalories)
	}
	if filter.MaxPrice > 0 {
		dbLocal = dbLocal.Where("price <= ?", filter.MaxPrice)
	}
	if len(allergens) > 0 {
		subQuery := db.Model(&models.Allergen{}).Where("LOWER(allergens.name) IN ?", allergens).Group("allergens.item_id").Select("COUNT(*) as num_allergens, allergens.item_id as item_id")
		dbLocal = dbLocal.Joins("FULL OUTER JOIN (?) as allergen_count ON allergen_count.item_id = menu_items.id", subQuery).Where("allergen_count.num_allergens = 0 OR allergen_count.num_allergens is NULL")
	}

	dbLocal.Find(&data)
	return data
}

func filterAllergens(allergens []string) []string {
	ret := []string{}
	for _, allergen := range allergens {
		if len(allergen) > 0 {
			ret = append(ret, allergen)
		}
	}
	return ret
}

// FetchItem retrieves an item with the given ID from the databse.
// Returns an error if the item can't be found.
func FetchItem(id int) (models.MenuItem, error) {
	ret := models.MenuItem{}
	res := db.Model(&models.MenuItem{}).Preload("Allergens").Where("ID = ?", id).First(&ret)
	if res.Error != nil {
		return models.MenuItem{}, res.Error
	}
	return ret, nil
}
