package api

import (
	"strings"

	"github.com/ZipCodeCore/zip-athena/constants"
	"github.com/ZipCodeCore/zip-athena/models"
)

func isJavaTopicCategory(category string) bool {
	if category == constants.TopicCategoryJava {
		return true
	}
	return strings.HasPrefix(category, constants.TopicCategoryJava+"-") ||
		strings.HasPrefix(category, constants.TopicCategoryJava+"_")
}

func isPythonTopicCategory(category string) bool {
	if category == constants.TopicCategoryPython {
		return true
	}
	return strings.HasPrefix(category, constants.TopicCategoryPython+"-") ||
		strings.HasPrefix(category, constants.TopicCategoryPython+"_")
}

func userCanAccessTopic(user *models.User, topic *models.Topic) bool {
	if user.IsAdmin {
		return true
	}

	switch user.Category {
	case constants.CategoryJava:
		return isJavaTopicCategory(topic.Category) || topic.Category == constants.TopicCategorySQL
	case constants.CategoryData:
		return isPythonTopicCategory(topic.Category) || topic.Category == constants.TopicCategorySQL
	default:
		return false
	}
}
