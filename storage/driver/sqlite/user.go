package sqlite

import "database/sql"

func (c *Config) AddUser(name string) {
	c.DB.Create(&User{Name: name})
}

func (c *Config) RemoveUser(name string) error {
	c.DB.Where("name = ?", name).Delete(User{})

	return nil
}

func (c *Config) ListUsers(filter map[string]string) map[string][]string {
	result := make(map[string][]string)

	q := c.DB.Table("users").Select("users.name, groups.name").Joins("LEFT JOIN group_users ON group_users.user_id = users.id").Joins("LEFT JOIN groups ON groups.id = group_users.group_id")

	if v, ok := filter["name"]; ok {
		q = q.Where("users.name = ?", v)
	}

	if v, ok := filter["elem"]; ok {
		q = q.Where("groups.name = ?", v)
	}

	rows, err := q.Rows()
	if err != nil {
		return result
	}
	defer rows.Close()

	for rows.Next() {
		var user string
		var group sql.NullString

		if err := rows.Scan(&user, &group); err != nil {
			continue
		}

		if _, ok := result[user]; !ok {
			result[user] = []string{}
		}
		if group.Valid {
			result[user] = append(result[user], group.String)
		}
	}

	return result
}

func (c *Config) FindUser(name string) bool {
	var count int64

	c.DB.Model(&User{}).Where("name = ?", name).Count(&count)

	return count == 1
}

func (c *Config) CountUser() int {
	var count int64

	c.DB.Model(&User{}).Count(&count)

	return int(count)
}

func (c *Config) AddUserToGroup(group, user string) {
	g := Group{}
	u := User{}

	c.DB.Where("name = ?", user).Find(&u)
	c.DB.Where("name = ?", group).Find(&g)

	c.DB.Model(&g).Association("Users").Append(&u)
}

func (c *Config) RemoveUserFromGroup(group, user string) {
	g := Group{}
	u := User{}

	c.DB.Where("name = ?", user).Find(&u)
	c.DB.Where("name = ?", group).Find(&g)

	c.DB.Model(&g).Association("Users").Delete(&u)
}
