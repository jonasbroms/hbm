package image

// From github.com/juliengk/go-docker/image

import (
	"fmt"
	"regexp"
	"strings"
)

type Image struct {
	ID       string
	Registry string
	Name     string
	Tag      string
	Official bool
}

func NewImage(img string) Image {
	name, tag := GetNameTag(img)

	image := Image{
		Name: name,
		Tag:  tag,
	}

	result := strings.Split(name, "/")
	count := len(result)

	if count >= 3 {
		image.Registry = result[0]
		image.Name = strings.Join(result[1:count], "/")
	} else if count == 2 {
		if validateRegistry(result[0]) {
			image.Registry = result[0]
			image.Name = result[1]
		}
	} else if count == 1 {
		image.Official = true
	}

	return image
}

func (img *Image) String() string {
	if img.Registry != "" && img.Name != "" {
		return fmt.Sprintf("%s/%s", img.Registry, img.Name)
	} else if img.Name != "" {
		return img.Name
	}

	return img.Name
}

func GetNameTag(name string) (string, string) {
	nt := strings.Split(name, ":")
	count := len(nt)

	if count > 2 {
		return strings.Join(nt[0:count-1], ":"), nt[count-1]
	} else if count == 2 {
		if strings.Contains(nt[1], "/") {
			return name, "latest"
		}
		return nt[0], nt[1]
	} else if count == 1 {
		return nt[0], "latest"
	}

	return "", ""
}

func validateRegistry(value string) bool {
	if len(value) == 0 {
		return false
	}

	result := strings.Split(value, ":")

	return isValidFQDN(result[0])
}

func isValidFQDN(s string) bool {
	if len(s) == 0 || len(s) > 254 {
		return false
	}
	parts := strings.Split(s, ".")
	if len(parts) < 2 {
		return false
	}
	reLabel := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9\-]{1,63}$`)
	reTLD := regexp.MustCompile(`^[a-zA-Z]{2,63}$`)
	for i, p := range parts {
		if i == len(parts)-1 {
			if !reTLD.MatchString(p) {
				return false
			}
		} else {
			if !reLabel.MatchString(p) {
				return false
			}
		}
	}
	return true
}
