package sqlite

import (
	"fmt"
)

func (c *Config) SetContainerOwner(username, name, containerid string) error {
	var user User
	c.DB.Where("name = ?", username).First(&user)
	if user.ID == 0 {
		return fmt.Errorf("user %q not found", username)
	}

	co := ContainerOwner{
		ContainerID:   containerid,
		ContainerName: name,
		User:          user,
	}
	return c.DB.Create(&co).Error
}

func (c *Config) IsContainerOwner(username, containerid string) bool {
	var u User
	c.DB.Where("name = ?", username).First(&u)
	if u.ID == 0 {
		return false
	}

	var count int

	// Match by container name (docker inspect/start/stop by name)
	c.DB.Model(&ContainerOwner{}).
		Where("container_name = ? AND user_id = ?", containerid, u.ID).
		Count(&count)
	if count > 0 {
		return true
	}

	// Match by full ID or short-ID prefix.
	// user_id is in the WHERE clause so a prefix shared with another user's
	// container never produces a false negative (fixes the old inverted LIKE bug).
	prefix := containerid + "%"
	c.DB.Model(&ContainerOwner{}).
		Where("container_id LIKE ? AND user_id = ?", prefix, u.ID).
		Count(&count)
	return count > 0
}

func (c *Config) RemoveContainerOwner(containerid string) error {
	return c.DB.Where("container_id = ?", containerid).Delete(&ContainerOwner{}).Error
}

func (c *Config) BackfillContainerName(containerID, name string) error {
	return c.DB.Model(&ContainerOwner{}).
		Where("container_id = ? AND (container_name IS NULL OR container_name = ?)", containerID, "").
		Update("container_name", name).Error
}

func (c *Config) ListContainerOwnerIDs() []string {
	var owners []ContainerOwner
	c.DB.Select("container_id").Find(&owners)
	ids := make([]string, 0, len(owners))
	for _, o := range owners {
		ids = append(ids, o.ContainerID)
	}
	return ids
}
