package attributes

// Taggable allows a struct to tag its origin, so it can be identified by other processes.
type Taggable interface {
	// GetTag returns a string that is used to tag itself, so it can be identified by other processes.
	GetTag() string
}
