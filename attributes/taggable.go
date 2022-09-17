package attributes

// Taggable allows a struct to tag its source, so it can be identified by other processes.
type Taggable interface {
	// GetOrigin returns a string that marks data as part of a group, so the source can be identified.
	GetOrigin() string
	// SetOrigin sets the origin which is used to identify its source.
	SetOrigin(origin string)
	// GetDataId returns a string that identifies where the data originates from, so it can be processed based on its source.
	GetDataId() string
	// SetDataId sets the source of the data so, it can be processed based on its dataId.
	SetDataId(dataId string)
}

// Tag implements Taggable by holding vars for its getters and setters.
type Tag struct {
	origin string
	dataId string
}

// NewTag returns a new instance of Tag.
func NewTag(origin string, dataId string) *Tag {
	return &Tag{
		origin,
		dataId,
	}
}

// GetOrigin implements Taggable.GetOrigin.
func (t *Tag) GetOrigin() string {
	return t.origin
}

// SetOrigin implements Taggable.SetOrigin.
func (t *Tag) SetOrigin(origin string) {
	t.origin = origin
}

// GetDataId implements Taggable.GetId.
func (t *Tag) GetDataId() string {
	return t.dataId
}

// SetDataId implements Taggable.SetId.
func (t *Tag) SetDataId(dataId string) {
	t.dataId = dataId
}
