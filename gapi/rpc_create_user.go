package gapi

import (
	"context"
	db "project/simplebank/db/sqlc"
	"project/simplebank/pb"
	"project/simplebank/util"
	"project/simplebank/val"
	"project/simplebank/worker"
	"time"

	"github.com/hibiken/asynq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	violations := validateCreateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	hashedPassword, err := util.HashedPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}
	arg := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			Username:       req.GetUsername(),
			FullName:       req.GetFullName(),
			Email:          req.GetEmail(),
			HashedPassword: hashedPassword,
		},
		AfterCreate: func(user db.User) error {
			taskPayload := &worker.PayloadSendVerifyEmail{
				Username: req.GetUsername(),
			}
			//重试次数为10 延迟处理
			options := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(time.Second * 5),
				asynq.Queue(worker.QueueCritical),
			}

			return server.taskDistributor.DistributeTaskSendVerifyEmail(ctx, taskPayload, options...)

		},
	}
	TxResult, err := server.store.CreateUserTx(ctx, arg)
	if err != nil {
		errCode := db.ErrorCode(err)
		//此处只保留一个外键约束
		if errCode == db.UniqueViolation {
			return nil, status.Errorf(codes.AlreadyExists, "username or email already exists: %s", err)

		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)

	}

	rsp := &pb.CreateUserResponse{
		User: convertUser(TxResult.User),
	}

	return rsp, nil
}

func validateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	if err := val.ValidateFullName(req.GetFullName()); err != nil {
		violations = append(violations, fieldViolation("full_name", err))
	}

	if err := val.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	return violations
}
