package database

const ScrapedDataTableName = "scraped_data"
const FilteredDataTableName = "filtered_data"

// Db is a facade that holds an instance of Driver and forwards its functions,
// Driver is interchangeable and allows the changing of database types.
// Db currently does not support context.Context.
type Db struct {
	Driver
}

// NewDb creates a new Db instance and attempts to connect the Driver to its database and
// Return an error in case connecting fails.
func NewDb(d Driver) (*Db, error) {
	err := d.connect()
	if err != nil {
		return nil, err
	}
	return &Db{
		d,
	}, err
}
