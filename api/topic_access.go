package api

import (
	"github.com/ZipCodeCore/zip-athena/constants"
	"github.com/ZipCodeCore/zip-athena/models"
)

func userCanAccessTopic(user *models.User, topic *models.Topic) bool {
	if user.IsAdmin {
		return true
	}

	if topic.Category == constants.TopicCategoryAll {
		return true
	}

	switch user.Category {
	case constants.CategoryJava:
		return topic.Category == constants.TopicCategoryJava
	case constants.CategoryData:
		return topic.Category == constants.TopicCategoryPython
	default:
		return false
	}
}
