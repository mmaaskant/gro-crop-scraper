// Set up the "scraped_data" table, which holds all raw data fetched during the crawler step.
db.scraped_data.drop()
db.scraped_data.createIndex({ url: 1 }, { unique: true })
db.scraped_data.createIndex({ config_id: 1 })
db.scraped_data.createIndex({ scraper_id: 1 })
db.scraped_data.createIndex({ created_at: 1 })
db.scraped_data.createIndex({ updated_at: 1 })

// Set up the "filtered_data" table, which holds the results of the filtered scraped data.
db.filtered_data.drop()
db.filtered_data.createIndex({ url: 1 }, { unique: true })
db.filtered_data.createIndex({ config_id: 1 })
db.filtered_data.createIndex({ scraper_id: 1 })
db.filtered_data.createIndex({ created_at: 1 })
db.filtered_data.createIndex({ updated_at: 1 })