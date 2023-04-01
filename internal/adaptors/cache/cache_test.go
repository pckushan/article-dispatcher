package cache

import (
	"article-dispatcher/internal/domain/adaptors/logger"
	"article-dispatcher/internal/domain/models"
	"article-dispatcher/internal/pkg/log"

	"github.com/stretchr/testify/assert"

	"context"
	"fmt"
	"sync"
	"testing"
)

func TestCache_Set(t *testing.T) {
	type fields struct {
		log             logger.Logger
		lock            *sync.RWMutex
		articles        map[string]models.Article
		tagDateIndexMap map[string][]string
	}
	type args struct {
		ctx     context.Context
		article *models.Article
	}
	l, err := log.NewLogger(log.ERROR)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	article := models.Article{
		Id:    "1",
		Title: "test",
		Date:  "2023-03-30",
		Body:  "test body",
		Tags:  []string{"fun", "health", "fitness"},
	}

	articlesCache := make(map[string]models.Article, 0)

	articlesCache[article.Id] = article

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "set_article_to_empty_cache",
			fields: fields{
				log:             l,
				lock:            &sync.RWMutex{},
				articles:        make(map[string]models.Article, 0),
				tagDateIndexMap: make(map[string][]string),
			},
			args: args{
				ctx:     context.Background(),
				article: &article,
			},
			wantErr: false,
		},
		{
			name: "set_article_to_id_already_existing_cache",
			fields: fields{
				log:             l,
				lock:            &sync.RWMutex{},
				articles:        articlesCache,
				tagDateIndexMap: make(map[string][]string),
			},
			args: args{
				ctx:     context.Background(),
				article: &article,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Cache{
				log:             tt.fields.log,
				lock:            tt.fields.lock,
				articles:        tt.fields.articles,
				tagDateIndexMap: tt.fields.tagDateIndexMap,
			}
			if err := c.Set(tt.args.ctx, tt.args.article); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCache_Get(t *testing.T) {
	type fields struct {
		log             logger.Logger
		lock            *sync.RWMutex
		articles        map[string]models.Article
		tagDateIndexMap map[string][]string
	}
	type args struct {
		ctx context.Context
		id  string
	}
	l, err := log.NewLogger(log.ERROR)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	article := models.Article{
		Id:    "1",
		Title: "test",
		Date:  "2023-03-30",
		Body:  "test body",
		Tags:  []string{"fun", "health", "fitness"},
	}

	articlesCache := make(map[string]models.Article, 0)
	tagDateIndexMap := make(map[string][]string)

	articlesCache[article.Id] = article
	tagDateIndexMap["fun#20230330"] = []string{"1"}
	tagDateIndexMap["health#20230330"] = []string{"1"}
	tagDateIndexMap["fitness#20230330"] = []string{"1"}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    models.Article
		wantErr bool
	}{
		{
			name: "get_article_from_cache",
			fields: fields{
				log:             l,
				lock:            &sync.RWMutex{},
				articles:        articlesCache,
				tagDateIndexMap: make(map[string][]string),
			},
			args: args{
				ctx: context.Background(),
				id:  "1",
			},
			wantErr: false,
			want:    article,
		},
		{
			name: "get_non_existing_article_from_cache",
			fields: fields{
				log:             l,
				lock:            &sync.RWMutex{},
				articles:        articlesCache,
				tagDateIndexMap: make(map[string][]string),
			},
			args: args{
				ctx: context.Background(),
				id:  "2",
			},
			wantErr: true,
			want:    models.Article{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Cache{
				log:             tt.fields.log,
				lock:            tt.fields.lock,
				articles:        tt.fields.articles,
				tagDateIndexMap: tt.fields.tagDateIndexMap,
			}
			got, err := c.Get(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// nolint: funlen
func TestCache_Filter(t *testing.T) {
	type fields struct {
		log             logger.Logger
		lock            *sync.RWMutex
		articles        map[string]models.Article
		tagDateIndexMap map[string][]string
	}
	type args struct {
		ctx  context.Context
		tag  string
		date int
	}

	l, err := log.NewLogger(log.ERROR)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	article := models.Article{
		Id:    "1",
		Title: "test",
		Date:  "2023-03-30",
		Body:  "test body",
		Tags:  []string{"fun", "health", "fitness"},
	}

	articlesCache := make(map[string]models.Article, 0)
	tagDateIndexMap := make(map[string][]string)

	articlesCache[article.Id] = article
	tagDateIndexMap["fun#20230330"] = []string{"1"}
	tagDateIndexMap["health#20230330"] = []string{"1"}
	tagDateIndexMap["fitness#20230330"] = []string{"1"}

	expectedTaggedArticle := models.TaggedArticles{
		Tag:         "health",
		Count:       3,
		Articles:    []string{"1"},
		RelatedTags: []string{"fun", "fitness"},
	}

	// cache maps to test more than 10 articals with the same tag
	articlesCacheMoreData := make(map[string]models.Article, 0)
	tagDateIndexMapMoreData := make(map[string][]string)

	for i := 1; i < 13; i++ {
		newID := fmt.Sprintf("%d", i)
		article.Id = newID
		articlesCacheMoreData[newID] = article

		if _, ok := tagDateIndexMapMoreData["fun#20230330"]; !ok {
			tagDateIndexMapMoreData["fun#20230330"] = []string{newID}
		} else {
			tagDateIndexMapMoreData["fun#20230330"] = append(tagDateIndexMapMoreData["fun#20230330"], newID)
		}

		if _, ok := tagDateIndexMapMoreData["health#20230330"]; !ok {
			tagDateIndexMapMoreData["health#20230330"] = []string{newID}
		} else {
			tagDateIndexMapMoreData["health#20230330"] = append(tagDateIndexMapMoreData["health#20230330"], newID)
		}
		if _, ok := tagDateIndexMapMoreData["fitness#20230330"]; !ok {
			tagDateIndexMapMoreData["fitness#20230330"] = []string{newID}
		} else {
			tagDateIndexMapMoreData["fitness#20230330"] = append(tagDateIndexMapMoreData["fitness#20230330"], newID)
		}
	}

	expectedTaggedArticleForMoreData := models.TaggedArticles{
		Tag:         "health",
		Count:       3,
		Articles:    []string{"3", "4", "5", "6", "7", "8", "9", "10", "11", "12"},
		RelatedTags: []string{"fun", "fitness"},
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    models.TaggedArticles
		wantErr bool
	}{
		{
			name: "filter_existing_article_from_cache",
			fields: fields{
				log:             l,
				lock:            &sync.RWMutex{},
				articles:        articlesCache,
				tagDateIndexMap: tagDateIndexMap,
			},
			args: args{
				ctx:  context.Background(),
				tag:  "health",
				date: 20230330,
			},
			wantErr: false,
			want:    expectedTaggedArticle,
		},
		{
			name: "filter_non_existing_article_from_cache",
			fields: fields{
				log:             l,
				lock:            &sync.RWMutex{},
				articles:        articlesCache,
				tagDateIndexMap: tagDateIndexMap,
			},
			args: args{
				ctx:  context.Background(),
				tag:  "invalid",
				date: 20230330,
			},
			wantErr: true,
			want: models.TaggedArticles{
				Articles:    make([]string, 0),
				RelatedTags: make([]string, 0),
			},
		},
		{
			name: "filter_tagged_article_with_more_than_limit",
			fields: fields{
				log:             l,
				lock:            &sync.RWMutex{},
				articles:        articlesCacheMoreData,
				tagDateIndexMap: tagDateIndexMapMoreData,
			},
			args: args{
				ctx:  context.Background(),
				tag:  "health",
				date: 20230330,
			},
			wantErr: false,
			want:    expectedTaggedArticleForMoreData,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Cache{
				log:             tt.fields.log,
				lock:            tt.fields.lock,
				articles:        tt.fields.articles,
				tagDateIndexMap: tt.fields.tagDateIndexMap,
			}
			got, err := c.Filter(tt.args.ctx, tt.args.tag, tt.args.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("Filter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !assert.Equal(t, tt.want, got) {
				t.Errorf("Filter() got = %v, want %v", got, tt.want)
			}
		})
	}
}
