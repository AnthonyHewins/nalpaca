package conf

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type Health struct {
	DisableHealth bool          `env:"DISABLE_HEALTH" envDefault:"false"`
	HealthPort    uint16        `env:"HEALTH_PORT" envDefault:"7654"`
	PingInterval  time.Duration `env:"HEALTH_PING_INTERVAL" envDefault:"10s"`
}

// This probably could just be a struct, the name never changes. No need to evaluate it
type HealthCheckable struct {
	Name string
	Fn   func(ctx context.Context) error
}

func (b *Bootstrapper) HealthServer(h *Health, dependencies ...HealthCheckable) *HealthServer {
	if h.DisableHealth {
		b.Logger.Info("health server disabled")
		return nil
	}

	s := &HealthServer{
		logger:        b.Logger,
		healthService: health.NewServer(),
		server:        grpc.NewServer(),
		port:          h.HealthPort,
		pingInterval:  h.PingInterval,
		dependencies:  dependencies,
	}

	healthpb.RegisterHealthServer(s.server, s.healthService)
	return s
}

type HealthServer struct {
	logger        *slog.Logger
	healthService *health.Server
	server        *grpc.Server
	port          uint16
	pingInterval  time.Duration
	dependencies  []HealthCheckable
}

func (s *HealthServer) GracefulStop() {
	s.server.GracefulStop()
}

// Start begins serving health
func (s *HealthServer) Start(ctx context.Context) error {
	addr := fmt.Sprintf(":%d", s.port)

	logger := s.logger.With("health_port", addr)

	// Reserve and attempt to listen on the specified port if its free.
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		logger.ErrorContext(ctx, "gRPC health server failed to listen on port", "error", err)
		return err
	}

	logger.InfoContext(ctx, "gRPC health server serving")

	// kill the goroutine to check health when the server goes down
	wrapper, cancel := context.WithCancel(ctx)
	defer cancel()

	go s.checkHealth(wrapper)
	return s.server.Serve(ln)
}

// checkHealth determines the health of a dependency.
func (s *HealthServer) checkHealth(ctx context.Context) {
	ticker := time.NewTicker(s.pingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.logger.Warn("health server thread ctx killed", "err", ctx.Err())
			return
		case <-ticker.C:
			for _, dep := range s.dependencies {
				name := dep.Name
				logger := s.logger.With("name", name)

				next := healthpb.HealthCheckResponse_SERVING
				if err := dep.Fn(ctx); err != nil {
					logger.ErrorContext(ctx, "dependency failed health check", "error", err)
					next = healthpb.HealthCheckResponse_NOT_SERVING
				}

				s.healthService.SetServingStatus(name, next)
			}
		}
	}
}
