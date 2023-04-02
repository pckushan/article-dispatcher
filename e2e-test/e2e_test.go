//go:build e2e

package e2e_test

import (
	cache2 "article-dispatcher/internal/adaptors/cache"
	"article-dispatcher/internal/domain/adaptors/logger"
	"article-dispatcher/internal/domain/models"
	"article-dispatcher/internal/http"
	"article-dispatcher/internal/http/responses"
	"article-dispatcher/internal/pkg/log"
	"article-dispatcher/internal/pkg/metrics"
	services2 "article-dispatcher/internal/services"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/stretchr/testify/assert"

	"context"
	"encoding/json"
	"fmt"
	netHttp "net/http"
	"strings"
	"testing"
)

func TestE2E(t *testing.T) {
	// load configurations and metrics
	initConfigs()
	// init new logger
	l, err := log.NewLogger(log.Config.Level)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// init and start metrics
	err = loadMetrics(l)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	cache := cache2.NewCache(l)

	// new fetcher and block services
	articleService := services2.NewArticleService(l, cache)

	// init router
	port := http.Config.Host
	router := http.Router{Conf: &http.Config}
	router.Init(l, articleService, metrics.RequestLatency)

	go func() {
		err := router.Start()
		t.Error(err)
	}()

	// sleep added to let the http server start
	time.Sleep(time.Second * 1)

	// router interrupt signal to stop router
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGTERM, os.Interrupt)
	go func() {
		<-s
		err := router.Stop()
		if err != nil {
			t.Error(err)
			return
		}
	}()

	// transaction request
	id := "1"
	tag := "health"
	date := 20230330
	runArticleCreateRequest(port, t)
	runArticleGetRequest(id, port, t)
	runArticleFilterRequest(tag, port, date, t)
}

func initConfigs() {
	// load logger configs
	log.Config = log.LoggerConfig{Level: "TRACE"}

	// load router configs
	http.Config = http.RouterConfig{Host: "8888"}
}

func loadMetrics(l logger.Logger) error {
	// load metrics configs
	mr := metrics.RouterMetrics{Conf: &metrics.Conf, Logger: l}
	err := mr.InitMetrics()
	if err != nil {
		return err
	}

	return nil
}

// do article create request
func runArticleCreateRequest(port string, t *testing.T) {
	// create url and request
	url := fmt.Sprintf("http://localhost:%s/articles", port)
	article := models.Article{
		Title: "test",
		Date:  "2023-03-30",
		Body:  "test body",
		Tags:  []string{"nature", "health"},
	}
	r, err := json.Marshal(article)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	reqBody := strings.NewReader(string(r))

	req, err := netHttp.NewRequestWithContext(context.Background(), netHttp.MethodPost, url, reqBody)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// expected response
	expectedResponse := responses.SuccessResponse{}
	expectedResponse.Data.ID = "1"

	client := netHttp.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer resp.Body.Close()

	// assert response code
	if !assert.Equal(t, netHttp.StatusCreated, resp.StatusCode) {
		t.Failed()
	}

	response := &responses.SuccessResponse{}
	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// assert response value
	if !assert.Equal(t, expectedResponse, *response) {
		t.Failed()
	}
}

// do article get request
func runArticleGetRequest(id, port string, t *testing.T) {
	// create url and request
	url := fmt.Sprintf("http://localhost:%s/articles/%s", port, id)

	req, err := netHttp.NewRequestWithContext(context.Background(), netHttp.MethodGet, url, netHttp.NoBody)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// expected response
	expectedArticle := models.Article{
		Id:    "1",
		Title: "test",
		Date:  "2023-03-30",
		Body:  "test body",
		Tags:  []string{"nature", "health"},
	}

	client := netHttp.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer resp.Body.Close()

	// assert response code
	if !assert.Equal(t, netHttp.StatusOK, resp.StatusCode) {
		t.Failed()
	}

	response := &models.Article{}
	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// assert response value
	if !assert.Equal(t, expectedArticle, *response) {
		t.Failed()
	}
}

// do article filter request
func runArticleFilterRequest(tag, port string, date int, t *testing.T) {
	// create url and request
	url := fmt.Sprintf("http://localhost:%s/tags/%s/%d", port, tag, date)

	req, err := netHttp.NewRequestWithContext(context.Background(), netHttp.MethodGet, url, netHttp.NoBody)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// expected response
	expectedTaggedArticle := models.TaggedArticles{
		Tag:         tag,
		Count:       2,
		Articles:    []string{"1"},
		RelatedTags: []string{"nature"},
	}

	client := netHttp.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer resp.Body.Close()

	// assert response code
	if !assert.Equal(t, netHttp.StatusOK, resp.StatusCode) {
		t.Failed()
	}

	response := &models.TaggedArticles{}
	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// assert response value
	if !assert.Equal(t, expectedTaggedArticle, *response) {
		t.Failed()
	}
}
