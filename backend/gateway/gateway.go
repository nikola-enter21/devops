package gateway

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/nikola-enter21/devops-fmi-course/api/gen/go/user/v1"
	"github.com/nikola-enter21/devops-fmi-course/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	log = logging.MustNewLogger()
)

func Serve(signalCtx context.Context, httpAddr, grpcTarget string) {
	gwCtx, gwCancel := context.WithCancel(context.Background())
	defer gwCancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := user.RegisterUserServiceHandlerFromEndpoint(gwCtx, mux, grpcTarget, opts); err != nil {
		log.Fatalf("failed to register gateway: %v", err)
	}

	httpServer := &http.Server{
		Addr:    httpAddr,
		Handler: allowCORS(mux),
	}

	go func() {
		log.Infow("HTTP gateway listening", "address", httpAddr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// wait for sigterm/sigint
	<-signalCtx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Infow("HTTP gateway shutting down...")
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Errorw("HTTP server shutdown error", "error", err)
	} else {
		log.Infow("HTTP gateway shutdown complete.")
	}
}

func allowCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		for _, o := range allowedOrigins() {
			if o == "*" || o == origin {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", strings.Join([]string{
			"GET", "POST", "PUT", "DELETE", "OPTIONS",
		}, ","))
		w.Header().Set("Access-Control-Allow-Headers", strings.Join([]string{
			"Origin", "Content-Type", "Accept", "Authorization",
		}, ","))
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			preflightHandler(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func preflightHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func allowedOrigins() []string {
	env := os.Getenv("ALLOWED_ORIGINS")
	if env == "" {
		return []string{"*"}
	}
	return strings.Split(env, ",")
}
