package models

// nolint:stylecheck
type Article struct {
	Id    string   `json:"id"`
	Title string   `json:"title"`
	Date  string   `json:"date" validate:"datetime=2006-01-02"`
	Body  string   `json:"body"`
	Tags  []string `json:"tags"`
}

type Articles []Article

type TaggedArticles struct {
	Tag         string   `json:"tag"`
	Count       int      `json:"count"`
	Articles    []string `json:"articles"`
	RelatedTags []string `json:"related_tags"`
}
