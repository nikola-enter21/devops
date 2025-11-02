package authorizer

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func UnaryServerInterceptor(auth Authorizer) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		rpcMethod := info.FullMethod

		role := extractRole(ctx)
		if role == "" {
			return nil, status.Error(codes.Unauthenticated, "missing role in metadata")
		}

		input := map[string]interface{}{
			"role": role,
			"rpc":  strings.TrimPrefix(rpcMethod, "/"),
		}

		allowed, err := auth.Evaluate(ctx, input)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "authorization error: %v", err)
		}
		if !allowed {
			return nil, status.Errorf(codes.PermissionDenied,
				"access denied for role %q on RPC %q", role, rpcMethod)
		}

		return handler(ctx, req)
	}
}

func StreamServerInterceptor(auth Authorizer) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		rpcMethod := info.FullMethod

		role := extractRole(ss.Context())
		if role == "" {
			return status.Error(codes.Unauthenticated, "missing role in metadata")
		}

		input := map[string]interface{}{
			"role": role,
			"rpc":  strings.TrimPrefix(rpcMethod, "/"),
		}

		allowed, err := auth.Evaluate(ss.Context(), input)
		if err != nil {
			return status.Errorf(codes.Internal, "authorization error: %v", err)
		}
		if !allowed {
			return status.Errorf(codes.PermissionDenied,
				"access denied for role %q on RPC %q", role, rpcMethod)
		}

		return handler(srv, ss)
	}
}

func extractRole(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	values := md.Get("x-user-role")
	if len(values) == 0 {
		return ""
	}
	return values[0]
}
