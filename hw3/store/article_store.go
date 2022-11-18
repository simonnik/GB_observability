package store

import (
	"context"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"go.uber.org/zap"

	"github.com/simonnik/GB_observability/hw3/e"
	"github.com/simonnik/GB_observability/hw3/l"
	"github.com/simonnik/GB_observability/hw3/m"
)

type ErrNotFound struct {
	Id string
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("article with id %s is not found", e.Id)
}

type ArticleStore struct {
	E      e.E
	logger *zap.Logger
	tracer opentracing.Tracer
}

func NewArticleStore(l *zap.Logger, t opentracing.Tracer) (ArticleStore, error) {
	e, err := e.NewE("articles", l, t)
	if err != nil {
		return ArticleStore{}, err
	}
	return ArticleStore{E: e, tracer: t}, nil
}

func (s ArticleStore) Add(c context.Context, article m.Article) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c, s.tracer,
		"store.Add")
	defer span.Finish()
	err := s.E.Insert(ctx, article)
	if err != nil {
		s.logger.Error(fmt.Sprintf(`store add error: %s`, err))
		span.LogFields(
			log.Error(err),
		)
		return err
	}
	span.LogFields(
		log.String("Success add article", fmt.Sprintf("%v", article)),
	)
	return nil
}

func (s ArticleStore) Search(c context.Context, query string) ([]m.Article, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c, s.tracer,
		"store.Search")
	defer span.Finish()
	result, err := s.E.Search(ctx, query)
	if err != nil {
		s.logger.Error(fmt.Sprintf(`store search error: %s`, err))
		span.LogFields(
			log.Error(err),
		)
		return nil, err
	}
	hits := result.Hits.Hits
	var articles []m.Article
	for _, hit := range hits {
		var article m.Article
		//map[string]interface{} -> struct
		err = mapstructure.Decode(hit.Source, &article)
		if err != nil {
			s.logger.Error(fmt.Sprintf(`cannot decode struct: %s`, err))
			span.LogFields(
				log.Error(err),
			)
			return nil, err
		}
		article.Id = hit.ID
		articles = append(articles, article)
	}
	span.LogFields(
		log.String("success query", query),
	)
	return articles, nil
}

func (s ArticleStore) Get(c context.Context, id string) (m.Article, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(c, s.tracer,
		"store.Get")
	defer span.Finish()
	result, err := s.E.Get(ctx, id)
	if err != nil {
		s.logger.Error(fmt.Sprintf(`store get error: %s`, err))
		span.LogFields(
			log.Error(err),
		)
		return m.Article{}, err
	}
	l.L(result)
	if found, ok := result["found"].(bool); ok {
		if !found {
			err = ErrNotFound{Id: id}
			//s.logger.Error(fmt.Sprintf(`article search error: %s`, err))
			span.LogFields(
				log.Error(err),
			)
			return m.Article{}, err
		}
	}

	var article m.Article
	err = mapstructure.Decode(result["_source"], &article)
	if err != nil {
		s.logger.Error(fmt.Sprintf(`cannot decode struct: %s`, err))
		span.LogFields(
			log.Error(err),
		)
		return m.Article{}, err
	}
	span.LogFields(
		log.String("success get article by id", id),
	)
	return article, nil
}
