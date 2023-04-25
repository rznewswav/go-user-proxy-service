package utils

import (
	"regexp"
)

func StringGetContentType(contentId string) string {
	if m, _ := regexp.MatchString("^A.*_.*$", contentId); m {
		// app_constants.ContentTypeArticle
		return "article"
	}
	if m, _ := regexp.MatchString("^V.*_.*$", contentId); m {
		// app_constants.ContentTypeVideo
		return "video"
	}
	if m, _ := regexp.MatchString("^P.*_.*$", contentId); m {
		// app_constants.ContentTypePodcast
		return "podcast"
	}
	if m, _ := regexp.MatchString("^F_.*_.*$", contentId); m {
		// app_constants.ContentTypeFeed
		return "feed"
	}
	// app_constants.ContentTypeArticle
	return "article"
}
