package conf

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type GrpcServerConf struct {
	EnableGrpc bool   `env:"ENABLE_GRPC" envDefault:"false"`
	GrpcPort   uint16 `env:"GRPC_PORT" envDefault:"9200"`

	MaxConnAge  time.Duration `env:"GRPC_MAX_CONN_AGE" envDefault:"6m"`
	MaxConnIdle time.Duration `env:"GRPC_MAX_CONN_IDLE" envDefault:"5m"`
	PingTimeout time.Duration `env:"GRPC_PING_TIMEOUT" envDefault:"0s"`
	PingEvery   time.Duration `env:"GRPC_PING_EVERY" envDefault:"2s"`
}

type GrpcServerConfWithProxy struct {
	GrpcServerConf

	EnableGrpcProxy bool   `env:"ENABLE_GRPC_PROXY" envDefault:"false"`
	GrpcProxyPort   uint16 `env:"GRPC_PROXY_PORT" envDefault:"8080"`

	MaxHeaderBytes         uint32        `env:"GRPC_PROXY_MAX_HEADER_BYTES"`
	IdleTimeout            time.Duration `env:"GRPC_PROXY_SERVER_TIMEOUT" envDefault:"120s"`
	ProxyReadHeaderTimeout time.Duration `env:"GRPC_PROXY_READ_HEADER_TIMEOUT" envDefault:"120s"`
	ProxyReadTimeout       time.Duration `env:"GRPC_PROXY_READ_TIMEOUT" envDefault:"120s"`
	ProxyWriteTimeout      time.Duration `env:"GRPC_PROXY_WRITE_TIMEOUT" envDefault:"120s"`
}

type GrpcServer struct {
	Server *grpc.Server
	Port   uint16
}

type GRPCGatewayHandler struct {
	Name    string                                                                                          // Name of the service used for logging purposes only.
	Handler func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error // Handler generated from protobuf file EX: foosvc.RegisterFooServiceHandlerFromEndpoint.
}

func (b *Bootstrapper) GRPC(ctx context.Context, g *GrpcServerConf, extraOpts ...grpc.ServerOption) GrpcServer {
	if !g.EnableGrpc {
		b.Logger.InfoContext(ctx, "grpc proxy disabled, skipping")
		return GrpcServer{}
	}

	l := b.Logger.With("config", g)

	extraOpts = append(extraOpts, grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: g.MaxConnIdle,
		MaxConnectionAge:  g.MaxConnAge,
		Time:              g.PingEvery,
		Timeout:           g.PingTimeout,
	}))

	l.Debug("creating grpc server")
	return GrpcServer{grpc.NewServer(extraOpts...), g.GrpcPort}
}

func (b *Bootstrapper) GrpcProxy(ctx context.Context, g *GrpcServerConfWithProxy, handlers ...GRPCGatewayHandler) (*http.Server, error) {
	if !g.EnableGrpcProxy || !g.EnableGrpc {
		b.Logger.InfoContext(ctx, "grpc proxy disabled either due to grpc being disabled or explicitly disabled")
		return nil, nil
	}

	l := b.Logger.With(
		"connectingPort", g.GrpcPort,
		"servingPort", g.GrpcProxyPort,
		"grpcProxyServerTimeout", g.IdleTimeout,
		"proxyReadHeaderTimeout", g.ProxyReadHeaderTimeout,
		"proxyReadTimeout", g.ProxyReadTimeout,
		"proxyWriteTimeout", g.ProxyWriteTimeout,
		"maxHeaderBytes", g.MaxHeaderBytes,
	)

	mux := runtime.NewServeMux()

	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	endpoint := fmt.Sprintf("localhost:%d", g.GrpcPort)
	for _, handler := range handlers {
		err := handler.Handler(ctx, mux, endpoint, dialOpts)
		if err != nil {
			l.ErrorContext(ctx, "failed binding "+handler.Name, "error", err)
			return nil, fmt.Errorf("registering '%s' service on grpc proxy: %w", handler.Name, err)
		}
	}

	l.InfoContext(ctx, "created grpc proxy server")
	return &http.Server{
		Addr:              fmt.Sprintf(":%d", g.GrpcProxyPort),
		Handler:           mux,
		ReadTimeout:       g.ProxyReadTimeout,
		ReadHeaderTimeout: g.ProxyReadHeaderTimeout,
		WriteTimeout:      g.ProxyWriteTimeout,
		IdleTimeout:       g.IdleTimeout,
		MaxHeaderBytes:    int(g.MaxHeaderBytes),
	}, nil
}
