package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/simonnik/GB_observability/hw3/m"
	"github.com/simonnik/GB_observability/hw3/store"
)

type ArticleHandler struct {
	S      store.ArticleStore
	logger *zap.Logger
	tracer opentracing.Tracer
}

func NewArticleHandler(s store.ArticleStore, l *zap.Logger, t opentracing.Tracer) ArticleHandler {
	return ArticleHandler{
		S:      s,
		logger: l,
		tracer: t,
	}
}

func (h ArticleHandler) Id(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c, h.tracer,
		"ArticleHandler.Get.Id")
	defer span.Finish()
	h.logger.Info("ArticleHandler.Get.Id called", zap.Field{Key: "method", String: c.Request.Method, Type: zapcore.StringType})
	span.SetTag("method", c.Request.Method)
	span.SetTag("params", c.Params)

	id := c.Params.ByName("id")

	article, err := h.S.Get(ctx, id)
	if err != nil {
		h.logger.Error(fmt.Sprintf(`failed to get article: %s`, err))
		span.LogFields(
			log.Error(err),
		)
		if errors.Is(err, store.ErrNotFound{Id: id}) {
			h.logger.Error(fmt.Sprintf(`article is not found: %s`, err))
			c.String(http.StatusNotFound, err.Error())
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	span.LogFields(
		log.String("Success get article", fmt.Sprintf("%v", article)),
	)
	c.JSON(http.StatusOK, article)
}

func (h ArticleHandler) Add(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(context.Background(), h.tracer,
		"ArticleHandler.Add")
	defer span.Finish()
	h.logger.Info("ArticleHandler.Post.Add called", zap.Field{Key: "method", String: c.Request.Method, Type: zapcore.StringType})
	span.SetTag("method", c.Request.Method)
	span.SetTag("params", c.Params)

	var article m.Article
	err := c.BindJSON(&article)
	if err != nil {
		h.logger.Error(fmt.Sprintf(`bad request: %s`, err))
		span.LogFields(
			log.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{"status": "bad request"})
		return
	} else {
		span.SetTag("body", article)
	}
	err = h.S.Add(ctx, article)
	if err != nil {
		h.logger.Error(fmt.Sprintf(`cannod add article: %s`, err))
		span.LogFields(
			log.Error(err),
		)
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	span.LogFields(
		log.String("Success add article", fmt.Sprintf("%v", article)),
	)
	c.JSON(http.StatusOK, article)
}

type SearchRequest struct {
	Query string `json:"query"`
}

func (h ArticleHandler) Search(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c, h.tracer,
		"ArticleHandler.Search")
	defer span.Finish()
	h.logger.Info("ArticleHandler.Post.Search called", zap.Field{Key: "method", String: c.Request.Method, Type: zapcore.StringType})
	span.SetTag("method", c.Request.Method)
	span.SetTag("params", c.Params)

	var query SearchRequest
	err := c.BindJSON(&query)
	if err != nil {
		h.logger.Error(fmt.Sprintf(`bad json: %s`, err))
		span.LogFields(
			log.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{"status": "bad request"})
		return
	}
	articles, err := h.S.Search(ctx, query.Query)
	if err != nil {
		h.logger.Error(fmt.Sprintf(`search error: %s`, err))
		span.LogFields(
			log.Error(err),
		)
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	span.LogFields(
		log.String("Success search articles", fmt.Sprintf("%v", articles)),
	)
	c.JSON(http.StatusOK, articles)
}
