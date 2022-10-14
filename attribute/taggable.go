package attribute

// Taggable allows a type to tag itself by storing a config.Config ID and a scraper.Scraper ID,
// these IDs allow processes to identify the source of the tagged type.
type Taggable interface {
	// GetConfigId returns the ID of the config.Config instance it was tagged by.
	GetConfigId() string
	// GetScraperId returns the id of the scraper.Scraper it was tagged by.
	GetScraperId() string
	// SetTag sets the Tag which is used to identify its source.
	SetTag(t *Tag)
}

// Tag partly implements Taggable.
type Tag struct {
	configId  string
	scraperId string
}

func NewTag(configId string, scraperId string) *Tag {
	return &Tag{
		configId,
		scraperId,
	}
}

// GetConfigId implements Taggable.GetConfigId.
func (t *Tag) GetConfigId() string {
	return t.configId
}

// GetScraperId implements Taggable.GetScraperId.
func (t *Tag) GetScraperId() string {
	return t.scraperId
}
