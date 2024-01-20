package main

import (
	"chs/config"
	configIn "chs/internal/config"
	General "chs/internal/general"
	"chs/internal/global_var"
	"chs/internal/init/subscription"
	"chs/internal/utils/cache"
	"chs/internal/utils/channel_manager"
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"testing"
	"time"

	"chs/pkg/discord"
	logger "chs/pkg/log"
	"chs/pkg/utils"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

var Version = "1.0.0"

var lock = &sync.Mutex{}

var serverInstance *http.Server
var ctx = context.Background()

func Debug(c *gin.Context) {

}

func TestTimeConsuming(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
}

func main() {
	// create root logger tagged with server version
	logX := logger.New().With(nil, "version", Version)
	configPath := utils.GetConfigPath(os.Getenv("config"))

	cfgFile, err := config.LoadConfig(configPath)

	if err != nil {
		logX.Fatalf("LoadConfig: %v", err)
	}
	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		logX.Fatalf("ParseConfig: %v", err)
	}
	logger := logger.InitOtelzap()
	global_var.AppLogger = logger

	global_var.PublicPath = cfg.Server.PublicPath
	global_var.Config = cfg
	// ActiveSubscription := subscription.ActiveSubscriptions
	if cfg.Jaeger.LogSpans {
		shutdown, err := InitOpenTelemetry(cfg)
		if err != nil {
			logger.Fatal("InitOpenTelemetry", zap.Error(err))
		}
		defer func() {
			if err := shutdown(ctx); err != nil {
				logger.Fatal("failed to shutdown TracerProvider", zap.Error(err))
			}
		}()
	}

	// Initialize Prometheus Metrics
	// metrics, err := metric.CreateMetrics(":9001", "pms-service")
	// if err != nil {
	// 	log.Fatal("Failed to create Prometheus metrics:", err)
	// }
	StartServer(cfg, logger)
	AwaitServerShutdown(cfg, logger)
}

func getInstance(config *config.Config, logger *otelzap.Logger) *http.Server {
	lock.Lock()
	defer lock.Unlock()
	if serverInstance == nil {
		logger.Info("serverInstance",
			zap.String("AppName", config.Server.AppName),
			zap.String("AppVersion", config.Server.AppVersion),
			zap.String("LogLevel", config.Logger.Level),
			zap.String("Mode", config.Server.Mode),
			zap.String("SSL", strconv.FormatBool(config.Server.SSL)))
		//load cache
		redisClient, err := configIn.NewRedisClient(ctx, config.Redis)
		if err != nil {
			logger.Fatal("NewRedisClient", zap.Error(err))
		}
		cache.InitCache(redisClient)
		channel_manager.CMInit(&config.CM)
		//load DB
		// DBMain, _, err := configIn.InitDB(config, logger)
		// if err != nil {
		// 	logger.Fatal("InitDB", zap.Error(err))
		// }

		// register the validation enum
		if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
			_ = v.RegisterValidation("enum", General.Enum)
		}

		// if !config.Server.Debug {
		subscription.LoadDataSubscription()
		// }
		serverInstance = &http.Server{
			// Handler: routes.HandleRequests(config.Server.Port, config.Server.Mode, config, DBMain, configIn.DBPool),
		}
	} else {
		logger.Debug("Server instance already exists")
	}

	return serverInstance
}

func StartServer(config *config.Config, logger *otelzap.Logger) {
	go func() {
		if err := getInstance(config, logger).ListenAndServe(); err != nil {
			logger.Debug(err.Error())
		}
	}()

	go discord.Init(config)
	// logger.Info("Server is listening at " + config.Server.Port)
}

func ShutdownServer(config *config.Config, logger *otelzap.Logger) {
	if serverInstance == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := getInstance(config, logger).Shutdown(ctx); err != nil {
		logger.Debug("Server forced to shutdown: " + err.Error())
	}
	logger.Info("Bye!")
}

func AwaitServerShutdown(config *config.Config, logger *otelzap.Logger) {
	quit := make(chan os.Signal, 10)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if config.Discord.Run {
		discord.UpdateChannelAPIStatus("\\🔴", " DOWN")
	}
	logger.Info("Shutting down server...")
	ShutdownServer(config, logger)
}

func InitOpenTelemetry(config *config.Config) (func(context.Context) error, error) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceName(config.Jaeger.ServiceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// If the OpenTelemetry Collector is running on a local cluster (minikube or
	// microk8s), it should be accessible through the NodePort service at the
	// `localhost:30080` endpoint. Otherwise, replace `localhost` with the
	// endpoint of your cluster. If you run the app inside k8s, then you can
	// probably connect directly to the service through dns.
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	fmt.Printf("Dialing Host: %s", config.Jaeger.Host)
	conn, err := grpc.DialContext(ctx, config.Jaeger.Host,
		// Note the use of insecure transport here. TLS is recommended in production.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	// Set up a trace exporter
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	sampler := sdktrace.ParentBased(sdktrace.NeverSample())
	if config.Jaeger.LogSpans {
		sampler = sdktrace.ParentBased(sdktrace.TraceIDRatioBased(config.Jaeger.SamplerRatio))
	}

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sampler),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Shutdown will flush any remaining spans and shut down the exporter.
	return tracerProvider.Shutdown, nil
}
