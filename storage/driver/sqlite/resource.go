package sqlite

import (
	"database/sql"

	"github.com/jonasbroms/hbm/object/types"
)

func (c *Config) AddResource(name, rtype, value, options string) {
	c.DB.Create(&Resource{
		Name:   name,
		Type:   rtype,
		Value:  value,
		Option: options,
	})
}

func (c *Config) RemoveResource(name string) error {
	c.DB.Where("name = ?", name).Delete(Resource{})

	return nil
}

func (c *Config) ListResources(filter map[string]string) map[types.Resource][]string {
	result := make(map[types.Resource][]string)

	q := c.DB.Table("resources").Select("resources.name, resources.type, resources.value, resources.option, collections.name").Joins("LEFT JOIN collection_resources ON collection_resources.resource_id = resources.id").Joins("LEFT JOIN collections ON collections.id = collection_resources.collection_id")

	if v, ok := filter["name"]; ok {
		q = q.Where("resources.name = ?", v)
	}

	if v, ok := filter["type"]; ok {
		q = q.Where("resources.type = ?", v)
	}

	if v, ok := filter["value"]; ok {
		q = q.Where("resources.value = ?", v)
	}

	if v, ok := filter["elem"]; ok {
		q = q.Where("collections.name = ?", v)
	}

	rows, _ := q.Rows()
	defer rows.Close()

	for rows.Next() {
		var resName string
		var resType string
		var resValue string
		var resOption string
		var collection sql.NullString

		if err := rows.Scan(&resName, &resType, &resValue, &resOption, &collection); err != nil {
			continue
		}

		rr := types.Resource{Name: resName, Type: resType, Value: resValue, Option: resOption}

		if _, ok := result[rr]; !ok {
			result[rr] = []string{}
		}
		if collection.Valid {
			result[rr] = append(result[rr], collection.String)
		}
	}

	return result
}

func (c *Config) FindResource(name string) bool {
	var count int64

	c.DB.Model(&Resource{}).Where("name = ?", name).Count(&count)

	return count == 1
}

func (c *Config) CountResource(rtype string) int {
	var count int64

	if rtype == "all" {
		c.DB.Model(&Resource{}).Count(&count)
	} else {
		c.DB.Model(&Resource{}).Where("type = ?", rtype).Count(&count)
	}

	return int(count)
}

func (c *Config) AddResourceToCollection(collection, resource string) {
	col := Collection{}
	res := Resource{}

	c.DB.Where("name = ?", resource).Find(&res)
	c.DB.Where("name = ?", collection).Find(&col)

	c.DB.Model(&col).Association("Resources").Append(&res)
}

func (c *Config) RemoveResourceFromCollection(collection, resource string) {
	col := Collection{}
	res := Resource{}

	c.DB.Where("name = ?", resource).Find(&res)
	c.DB.Where("name = ?", collection).Find(&col)

	c.DB.Model(&col).Association("Resources").Delete(&res)
}
