package gapi

import (
	"context"

	db "github.com/JeongWoo-Seo/simpleBank/db/sqlc"
	"github.com/JeongWoo-Seo/simpleBank/pb"
	"github.com/JeongWoo-Seo/simpleBank/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	violation := validDataVerifyEmailRequest(req)
	if violation != nil {
		return nil, invalidArgumentError(violation)
	}

	txResult, err := s.store.VerifyEmailTx(ctx, db.VerifyEmailTxParams{
		EmailID:    req.GetEmailId(),
		SecretCode: req.GetSecretCode(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to verify email")
	}

	res := &pb.VerifyEmailResponse{
		IsVerify: txResult.User.IsEmailVerified,
	}

	return res, nil
}

func validDataVerifyEmailRequest(req *pb.VerifyEmailRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidDataEmailID(req.GetEmailId()); err != nil {
		violations = append(violations, fieldViolation("email_id", err))
	}

	if err := val.ValidDataSecretCode(req.GetSecretCode()); err != nil {
		violations = append(violations, fieldViolation("secret_code", err))
	}

	return violations
}
