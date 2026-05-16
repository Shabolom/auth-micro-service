package di

import (
	authv1 "auth-micro-service/gen"
	"auth-micro-service/pkg/utils"
	"context"
	"errors"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

func (d *DI) NewAuthGRPCServer(logger *zap.Logger, authHandlers authv1.AuthServiceServer, usersHandlers authv1.UserServiceServer) *grpc.Server {
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(d.loggingInterceptor(logger)),
	)

	authv1.RegisterAuthServiceServer(grpcServer, authHandlers)
	authv1.RegisterUserServiceServer(grpcServer, usersHandlers)

	reflection.Register(grpcServer)

	logger.Info("grpc server initialized")

	return grpcServer
}

func (d *DI) loggingInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		start := time.Now()

		switch info.FullMethod {
		case authv1.AuthService_Register_FullMethodName,
			authv1.AuthService_Refresh_FullMethodName,
			authv1.AuthService_Login_FullMethodName:
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			logger.Info("failed to extract metadata")
			return nil, status.Error(codes.Unauthenticated, "metadata missing")
		}

		authorization := md.Get("authorization")

		if len(authorization) == 0 {
			logger.Info("authorization missing", zap.String("authorization", info.FullMethod))
			return nil, status.Error(codes.Unauthenticated, "token missing")
		}

		token := strings.TrimPrefix(authorization[0], "Bearer ")

		claims, err := utils.ParseToken(token, d.Config().Secret, d.Logger())
		if err != nil {
			logger.Info("failed to parse token", zap.String("token", token), zap.Error(err))
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		err = d.GetRedisHandlers().CheckSessionStatus(d.ctx, claims.ID)
		if err == redis.Nil || err != nil {
			d.Logger().Info("check error :", zap.String("id", claims.ID), zap.Error(err))
			return nil, status.Error(codes.Unauthenticated, errors.New("user is not authorized").Error())
		}

		resp, err := handler(ctx, req)
		if err != nil {
			logger.Info("handler returned error", zap.Error(err))
			return nil, status.Error(codes.Internal, err.Error())
		}

		logger.Info("grpc request",
			zap.String("method", info.FullMethod),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)

		return resp, err
	}
}
