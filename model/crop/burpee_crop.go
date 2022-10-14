package crop

import (
	"github.com/mmaaskant/gro-crop-scraper/database"
)

type BurpeeCrop struct {
	*database.Entity
}

func (bc *BurpeeCrop) GetCategory() *string {
	return bc.GetString("category")
}

func (bc *BurpeeCrop) GetFamily() *string {
	return bc.GetString("") // TODO: Split family and name
}

func (bc *BurpeeCrop) GetName() *string {
	return bc.GetString("") // TODO: Split name and family
}

func (bc *BurpeeCrop) GetLatinName() *string {
	return bc.GetString("latin_name")
}

func (bc *BurpeeCrop) GetDescription() *string {
	return bc.GetString("description")
}

func (bc *BurpeeCrop) GetShortDescription() *string {
	return bc.GetString("short_description")
}

func (bc *BurpeeCrop) GetGrowingZoneRange() *string {
	return bc.GetString("") // TODO: Handle slice
}

func (bc *BurpeeCrop) GetSunRequirement() *string {
	return bc.GetString("") // TODO: Handle slice
}

func (bc *BurpeeCrop) GetDaysToMaturity() *string {
	return bc.GetString("bp_days_to_maturity")
}

func (bc *BurpeeCrop) GetFruitSize() *string {
	return bc.GetString("bp_fruit_size")
}

func (bc *BurpeeCrop) GetMatureHeight() *string {
	return bc.GetString("bp_mature_height")
}

func (bc *BurpeeCrop) GetMatureSpread() *string {
	return bc.GetString("bp_mature_spread")
}
