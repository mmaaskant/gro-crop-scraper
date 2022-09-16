package attributes

// Taggable allows a struct to tag its origin, so it can be identified by other processes.
type Taggable interface {
	// GetTag returns a string that marks data as part of a group, so the source can be identified.
	GetTag() string
	// SetTag sets the Taggable tag which is used to identify its source.
	SetTag(tag string)
	// GetOrigin returns a string that identifies where the data originates from, so it can be processed based on its source.
	GetOrigin() string
}
