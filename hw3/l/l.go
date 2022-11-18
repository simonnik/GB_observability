package l

import (
	"encoding/json"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
	"io"
	"log"
	"reflect"
)

type zapWrapper struct {
	logger *zap.Logger
}

// Error logs a message at error priority
func (w *zapWrapper) Error(msg string) {
	w.logger.Error(msg)
}

// Infof logs a message at info priority
func (w *zapWrapper) Infof(msg string, args ...interface{}) {
	w.logger.Sugar().Infof(msg, args...)
}

// InitJaeger Для инициализации через передаваемые параметры
//func InitJaeger(service string, jaegerServerHostPort string, logger *zap.Logger) (opentracing.Tracer, io.Closer) {
//	cfg := &config.Configuration{
//		ServiceName: service,
//		Sampler: &config.SamplerConfig{
//			Type:  "const",
//			Param: 1,
//		},
//		Reporter: &config.ReporterConfig{
//			LogSpans:           true,
//			LocalAgentHostPort: jaegerServerHostPort,
//		},
//	}
//
//	tracer, closer, err := cfg.NewTracer(config.Logger(&zapWrapper{logger: logger}))
//	if err != nil {
//		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
//	}
//	return tracer, closer
//}

// InitJaeger Для инициализации через переменные окружения
func InitJaeger(logger *zap.Logger) (opentracing.Tracer, io.Closer) {
	cfg, err := config.FromEnv()
	if err != nil {
		logger.Warn("cannot parse Jaeger env vars", zap.Error(err))
	}
	tracer, closer, err := cfg.NewTracer(config.Logger(&zapWrapper{logger: logger}))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	return tracer, closer
}

var LogFunc = log.Println
var FatalFunc = log.Fatal

func F(i interface{}) {
	FatalFunc(Parse(i))
}
func Parse(i interface{}) interface{} {
	typeOf := reflect.TypeOf(i)
	if typeOf.Kind() == reflect.Map {
		b, err := json.MarshalIndent(i, "", "  ")
		if err != nil {
			return i
		}
		return string(b)
	}
	if typeOf.Kind() == reflect.Struct {
		b, err := json.MarshalIndent(i, "", "  ")
		if err != nil {
			return i
		}
		name := typeOf.Name()
		result := fmt.Sprintf("<%v>%v", name, string(b))
		return result
	}
	if typeOf.Kind() == reflect.Slice {
		v := reflect.ValueOf(i)
		result := "["
		for i := 0; i < v.Len(); i++ {
			val := v.Index(i)
			result += fmt.Sprintf("\n%v,", Parse(val.Interface()))
		}
		result = result[:len(result)-1] + "\n"
		result += "]"
		return result
	}
	if typeOf.Kind() == reflect.Array {
		v := reflect.ValueOf(i)
		result := "["
		for i := 0; i < v.Len(); i++ {
			val := v.Index(i)
			result += fmt.Sprintf("%v,", Parse(val.Interface()))
		}
		result += "]"
		return result
	}
	return i
}
func L(i interface{}) {
	LogFunc(Parse(i))
}
