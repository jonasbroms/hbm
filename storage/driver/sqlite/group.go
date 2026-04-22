package sqlite

import "database/sql"

func (c *Config) AddGroup(name string) {
	c.DB.Create(&Group{Name: name})
}

func (c *Config) RemoveGroup(name string) error {
	c.DB.Where("name = ?", name).Delete(Group{})

	return nil
}

func (c *Config) ListGroups(filter map[string]string) map[string][]string {
	result := make(map[string][]string)

	q := c.DB.Table("groups").Select("groups.name, users.name").Joins("LEFT JOIN group_users ON group_users.group_id = groups.id").Joins("LEFT JOIN users ON users.id = group_users.user_id")

	if v, ok := filter["name"]; ok {
		q = q.Where("groups.name = ?", v)
	}

	if v, ok := filter["elem"]; ok {
		q = q.Where("users.name = ?", v)
	}

	rows, _ := q.Rows()
	defer rows.Close()

	for rows.Next() {
		var group string
		var user sql.NullString

		if err := rows.Scan(&group, &user); err != nil {
			continue
		}

		if _, ok := result[group]; !ok {
			result[group] = []string{}
		}
		if user.Valid {
			result[group] = append(result[group], user.String)
		}
	}

	return result
}

func (c *Config) FindGroup(name string) bool {
	var count int64

	c.DB.Model(&Group{}).Where("name = ?", name).Count(&count)

	return count == 1
}

func (c *Config) CountGroup() int {
	var count int64

	c.DB.Model(&Group{}).Count(&count)

	return int(count)
}
