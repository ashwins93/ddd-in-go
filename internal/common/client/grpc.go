package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"os"
	"strconv"
	"time"

	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/genproto/trainer"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/genproto/users"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func NewTrainerClient() (client trainer.TrainerServiceClient, close func() error, err error) {
	grpcAddr := os.Getenv("TRAINER_GRPC_ADDR")
	if grpcAddr == "" {
		return nil, func() error { return nil }, errors.New("empty env TRAINER_GRPC_ADDR")
	}

	opts, err := grpcDialOpts(grpcAddr)
	if err != nil {
		return nil, func() error { return nil }, err
	}

	conn, err := grpc.Dial(grpcAddr, opts...)
	if err != nil {
		return nil, func() error { return nil }, err
	}

	return trainer.NewTrainerServiceClient(conn), conn.Close, nil
}

func WaitForTrainerService(timeout time.Duration) bool {
	return grpcHealthCheck(os.Getenv("TRAINER_GRPC_ADDR"), timeout)
}

func NewUsersClient() (client users.UsersServiceClient, close func() error, err error) {
	grpcAddr := os.Getenv("USERS_GRPC_ADDR")
	if grpcAddr == "" {
		return nil, func() error { return nil }, errors.New("empty env USERS_GRPC_ADDR")
	}

	opts, err := grpcDialOpts(grpcAddr)
	if err != nil {
		return nil, func() error { return nil }, err
	}

	conn, err := grpc.Dial(grpcAddr, opts...)
	if err != nil {
		return nil, func() error { return nil }, err
	}

	return users.NewUsersServiceClient(conn), conn.Close, nil
}

func WaitForUsersService(timeout time.Duration) bool {
	return grpcHealthCheck(os.Getenv("USERS_GRPC_ADDR"), timeout)
}

func grpcDialOpts(grpcAddr string) ([]grpc.DialOption, error) {
	if noTLS, _ := strconv.ParseBool(os.Getenv("GRPC_NO_TLS")); noTLS {
		return []grpc.DialOption{grpc.WithInsecure()}, nil
	}

	systemRoots, err := x509.SystemCertPool()
	if err != nil {
		return nil, errors.Wrap(err, "cannot load root CA cert")
	}
	creds := credentials.NewTLS(&tls.Config{
		RootCAs:    systemRoots,
		MinVersion: tls.VersionTLS12,
	})

	return []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(newMetadataServerToken(grpcAddr)),
	}, nil
}

func grpcHealthCheck(grpcAddr string, timeout time.Duration) bool {
	logrus.Debugf("checking if service %q is ready", grpcAddr)
	opts, err := grpcDialOpts(grpcAddr)
	if err != nil {
		logrus.Errorf("error: failed to create grpc dial options: %+v", err)
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	conn, err := grpc.DialContext(ctx, grpcAddr, opts...)
	if err != nil {
		if err == context.DeadlineExceeded {
			logrus.Errorf("timeout: failed to connect service %q within %v", grpcAddr, timeout)
		} else {
			logrus.Errorf("error: failed to connect service at %q: %+v", grpcAddr, err)
		}
		return false
	}
	defer conn.Close()
	logrus.Infof("service %q is ready", grpcAddr)

	return true
}
