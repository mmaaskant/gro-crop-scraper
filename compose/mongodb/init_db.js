// Set up the "scraped_data" table, which holds all raw data fetched during the crawler step.
db.scraped_data.drop()
db.scraped_data.createIndex({ url: 1 }, { unique: true })
db.scraped_data.createIndex({ config_id: 1 })
db.scraped_data.createIndex({ scraper_id: 1 })
db.scraped_data.createIndex({ data_type: 1 })
db.scraped_data.createIndex({ created_at: 1 })
db.scraped_data.createIndex({ updated_at: 1 })

// TODO: INIT filtered_data table