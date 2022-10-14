package model

import "github.com/mmaaskant/gro-crop-scraper/database"

type Crop interface {
	GetCategory() *string
	GetFamily() *string
	GetName() *string
	GetLatinName() *string
	GetDescription() *string
	GetShortDescription() *string
	GetGrowingZoneRange() *string
	GetSunRequirement() *string
	GetDaysToMaturity() *string
	GetFruitSize() *string
	GetMatureHeight() *string
	GetMatureSpread() *string
}

type GroCrop struct {
	*database.Entity
}

func NewGroCrop(e *database.Entity) *GroCrop {
	return &GroCrop{
		e,
	}
}

func (gc *GroCrop) GetCategory() *string {
	return gc.GetString("category")
}

func (gc *GroCrop) GetFamily() *string {
	return gc.GetString("family")
}

func (gc *GroCrop) GetName() *string {
	return gc.GetString("name")
}

func (gc *GroCrop) GetLatinName() *string {
	return gc.GetString("latin_name")
}

func (gc *GroCrop) GetDescription() *string {
	return gc.GetString("description")
}

func (gc *GroCrop) GetShortDescription() *string {
	return gc.GetString("short_description")
}

func (gc *GroCrop) GetGrowingZoneRange() *string {
	return gc.GetString("growing_zone_range")
}

func (gc *GroCrop) GetSunRequirement() *string {
	return gc.GetString("sun_requirement")
}

func (gc *GroCrop) GetDaysToMaturity() *string {
	return gc.GetString("days_to_maturity")
}

func (gc *GroCrop) GetFruitSize() *string {
	return gc.GetString("fruit_size")
}

func (gc *GroCrop) GetMatureHeight() *string {
	return gc.GetString("mature_height")
}

func (gc *GroCrop) GetMatureSpread() *string {
	return gc.GetString("mature_spread")
}
