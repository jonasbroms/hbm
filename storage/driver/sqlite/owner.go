package sqlite

import (
	"fmt"
)

func (c *Config) SetContainerOwner(username, name, containerid string) error {
	var user User

	c.DB.Where("name = ?", username).First(&user)
	if user.ID == 0 {
		return nil //FIXME
	}

	co := ContainerOwner{
		ContainerID: containerid,
		User:        user,
	}

	if len(name) > 1 { // container name must be at least 2 characters
		co.Name = name
	}
	c.DB.Model(&ContainerOwner{}).Create(&co)

	return nil
}

func (c *Config) IsContainerOwner(username, containerid string) bool {
	var co ContainerOwner
	var u User
	var cnt int

	c.DB.Where("name = ?", username).First(&u)
	if u.ID == 0 {
		return false
	}

	c.DB.Model(&co).Where("container_id = ? AND user_id = ?", containerid, u.ID).Count(&cnt)
	if cnt == 1 {
		return true
	}

	c.DB.Model(&co).Where("name = ? AND user_id = ?", containerid, u.ID).Count(&cnt)
	if cnt == 1 {
		return true
	}

	prefix := fmt.Sprintf("%s%%", containerid)
	prfm := false
	var cop []ContainerOwner
	c.DB.Where("container_id LIKE ?", prefix).Find(&cop)

	for _, p := range cop {
		if p.UserID != u.ID {
			return false
		}
		prfm = true
	}

	return prfm
}
