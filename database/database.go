package database

// Db is a facade that holds an instance of Driver and forwards its functions.
// Driver is interchangeable and allows the changing of database types.
type Db struct {
	Driver
}

// NewDb creates a new Db instance and attempts to connect the Driver to its database.
// Return an error in case connecting fails.
func NewDb(d Driver) (*Db, error) {
	err := d.Connect()
	if err != nil {
		return nil, err
	}
	return &Db{
		d,
	}, err
}
