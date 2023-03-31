package models

type Article struct {
	Id    string   `json:"id"`
	Title string   `json:"title"`
	Date  string   `json:"date"`
	Body  string   `json:"body"`
	Tags  []string `json:"tags"`
}

type Articles []Article

func NewArticleList() Articles {
	return make([]Article, 0)
}
