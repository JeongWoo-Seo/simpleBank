package gapi

import (
	"context"
	"errors"

	db "github.com/JeongWoo-Seo/simpleBank/db/sqlc"
	"github.com/JeongWoo-Seo/simpleBank/pb"
	"github.com/JeongWoo-Seo/simpleBank/util"
	"github.com/JeongWoo-Seo/simpleBank/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	authPayload, err := s.authorizationUser(ctx, []string{util.DepositorRole, util.BankerRole})
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	violation := validDataCreateAccountRequest(req)
	if violation != nil {
		return nil, invalidArgumentError(violation)
	}

	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Balance:  0,
		Currency: req.GetCurrency(),
	}

	account, err := s.store.CreateAccount(ctx, arg)
	if err != nil {
		if db.ErrorCode(err) == db.ForeignkeyViolation {
			return nil, status.Errorf(codes.NotFound, "user does not exist")
		}
		return nil, status.Errorf(codes.Internal, "fail create account : %s ", err)
	}

	res := &pb.CreateAccountResponse{
		Account: converterAccount(account),
	}

	return res, nil
}

func (s *Server) GetAccount(ctx context.Context, req *pb.GetAccountRequest) (*pb.GetAccountResponse, error) {
	authPayload, err := s.authorizationUser(ctx, []string{util.DepositorRole, util.BankerRole})
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	violation := validDataGetAccountRequest(req)
	if violation != nil {
		return nil, invalidArgumentError(violation)
	}

	account, err := s.store.GetAccount(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "account no exist")
		}
		return nil, status.Errorf(codes.Internal, "failed to get account : %s ", err)
	}

	if authPayload.Role != util.BankerRole && authPayload.Username != account.Owner {
		return nil, status.Errorf(codes.PermissionDenied, "you don't have permission to view this account")
	}

	res := &pb.GetAccountResponse{
		Account: converterAccount(account),
	}

	return res, nil
}

func (s *Server) ListAccount(ctx context.Context, req *pb.ListAccountRequest) (*pb.ListAccountResponse, error) {
	authPayload, err := s.authorizationUser(ctx, []string{util.DepositorRole, util.BankerRole})
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	violation := validDataListAccountRequest(req)
	if violation != nil {
		return nil, invalidArgumentError(violation)
	}

	arg := db.ListAccountsParams{
		Owner:  authPayload.Username,
		Limit:  req.GetPageSize(),
		Offset: (req.GetPageId() - 1) * req.GetPageSize(),
	}

	accountList, err := s.store.ListAccounts(ctx, arg)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "account no exist")
		}
		return nil, status.Errorf(codes.Internal, "failed to get account : %s ", err)
	}

	var pbAccount []*pb.Account
	for _, account := range accountList {
		pbAccount = append(pbAccount, converterAccount(account))
	}
	res := &pb.ListAccountResponse{
		Account: pbAccount,
	}

	return res, nil
}

func validDataCreateAccountRequest(req *pb.CreateAccountRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidDateCurrency(req.GetCurrency()); err != nil {
		violations = append(violations, fieldViolation("currency", err))
	}

	return violations
}

func validDataGetAccountRequest(req *pb.GetAccountRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidDataID(req.GetId()); err != nil {
		violations = append(violations, fieldViolation("ID", err))
	}

	return violations
}

func validDataListAccountRequest(req *pb.ListAccountRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidDatePageID(req.GetPageId()); err != nil {
		violations = append(violations, fieldViolation("page_id", err))
	}

	if err := val.ValidDatePageSize(req.GetPageSize()); err != nil {
		violations = append(violations, fieldViolation("page_size", err))
	}

	return violations
}
