package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/simonnik/GB_observability/hw2/m"
	"github.com/simonnik/GB_observability/hw2/store"
)

type ArticleHandler struct {
	S store.ArticleStore
}

func NewArticleHandler(s store.ArticleStore) ArticleHandler {
	return ArticleHandler{S: s}
}

func (h ArticleHandler) Id(c *gin.Context) {
	id := c.Params.ByName("id")
	ctx := context.Background()
	article, err := h.S.Get(ctx, id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, article)
}

func (h ArticleHandler) Add(c *gin.Context) {
	ctx := context.Background()
	var article m.Article
	if c.BindJSON(&article) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "bad request"})
		return
	}
	err := h.S.Add(ctx, article)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, article)
}

type SearchRequest struct {
	Query string `json:"query"`
}

func (h ArticleHandler) Search(c *gin.Context) {
	ctx := context.Background()
	var query SearchRequest
	if c.BindJSON(&query) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "bad request"})
		return
	}
	articles, err := h.S.Search(ctx, query.Query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, articles)
}
