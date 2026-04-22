package sqlite

import "database/sql"

func (c *Config) AddCollection(name string) {
	c.DB.Create(&Collection{Name: name})
}

func (c *Config) RemoveCollection(name string) error {
	c.DB.Where("name = ?", name).Delete(Collection{})

	return nil
}

func (c *Config) ListCollections(filter map[string]string) map[string][]string {
	result := make(map[string][]string)

	q := c.DB.Table("collections").Select("collections.name, resources.name").Joins("LEFT JOIN collection_resources ON collection_resources.collection_id = collections.id").Joins("LEFT JOIN resources ON resources.id = collection_resources.resource_id")

	if v, ok := filter["name"]; ok {
		q = q.Where("collections.name = ?", v)
	}

	if v, ok := filter["elem"]; ok {
		q = q.Where("resources.name = ?", v)
	}

	rows, _ := q.Rows()
	defer rows.Close()

	for rows.Next() {
		var collection string
		var resource sql.NullString

		if err := rows.Scan(&collection, &resource); err != nil {
			continue
		}

		if _, ok := result[collection]; !ok {
			result[collection] = []string{}
		}
		if resource.Valid {
			result[collection] = append(result[collection], resource.String)
		}
	}

	return result
}

func (c *Config) FindCollection(name string) bool {
	var count int64

	c.DB.Model(&Collection{}).Where("name = ?", name).Count(&count)

	return count == 1
}

func (c *Config) CountCollection() int {
	var count int64

	c.DB.Model(&Collection{}).Count(&count)

	return int(count)
}
