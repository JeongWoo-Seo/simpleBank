package gapi

import (
	"context"
	"errors"
	"fmt"
	"time"

	db "github.com/JeongWoo-Seo/simpleBank/db/sqlc"
	"github.com/JeongWoo-Seo/simpleBank/pb"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) RenewAccessToken(ctx context.Context, req *pb.RenewAccessTokenRequest) (*pb.RenewAccessTokenResponse, error) {

	violation := validDataRenewAccessTokenRequest(req)
	if violation != nil {
		return nil, invalidArgumentError(violation)
	}

	refreshPayload, err := s.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	session, err := s.store.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "not found session: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to get session: %s", err)
	}

	if session.IsBlocked {
		err := fmt.Errorf("blocked session")
		return nil, unauthenticatedError(err)
	}

	if session.Username != refreshPayload.Username {
		err := fmt.Errorf("incorrect session user")
		return nil, unauthenticatedError(err)
	}

	if session.RefreshToken != req.RefreshToken {
		err := fmt.Errorf("incorrect refresh token")
		return nil, unauthenticatedError(err)
	}

	if time.Now().After(session.ExpiresAt) {
		err := fmt.Errorf("expired session")
		return nil, unauthenticatedError(err)
	}

	accessToken, accessPayload, err := s.tokenMaker.CreateToken(refreshPayload.Username, refreshPayload.Role, s.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create accessToken: %s", err)
	}

	res := &pb.RenewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: timestamppb.New(accessPayload.ExpiredAt),
	}

	return res, nil
}

func validDataRenewAccessTokenRequest(req *pb.RenewAccessTokenRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if len(req.GetRefreshToken()) <= 0 {
		err := fmt.Errorf("refresh token is required and cannot be empty")
		violations = append(violations, fieldViolation("refreshToken", err))
	}

	return violations
}
