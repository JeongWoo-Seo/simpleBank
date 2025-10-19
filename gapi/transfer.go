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

func (s *Server) CreateTransfer(ctx context.Context, req *pb.CreateTransferRequest) (*pb.CreateTransferResponse, error) {
	authPayload, err := s.authorizationUser(ctx, []string{util.DepositorRole, util.BankerRole})
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	violation := validDataTransferRequest(req)
	if violation != nil {
		return nil, invalidArgumentError(violation)
	}

	fromAcouunt, err := s.validAccount(ctx, req.GetFromAccountId(), req.GetCurrency())
	if err != nil {
		return nil, err
	}

	if fromAcouunt.Owner != authPayload.Username {
		err := errors.New("from account doent belong to the authenticated user")
		return nil, unauthenticatedError(err)
	}
	_, err = s.validAccount(ctx, req.GetToAccountId(), req.GetCurrency())
	if err != nil {
		return nil, err
	}

	arg := db.TransferTxParams{
		FromAccountID: req.GetFromAccountId(),
		ToAccountID:   req.GetToAccountId(),
		Amount:        req.GetAmount(),
	}

	result, err := s.store.TransferTx(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to transfer tx: %s", err)
	}

	res := &pb.CreateTransferResponse{
		Transfer:    converterTransfer(result.Transfer),
		FromAccount: converterAccount(result.FromAccount),
		ToAccount:   converterAccount(result.ToAccount),
		FromEntry:   converterEntry(result.FromEntry),
		ToEntry:     converterEntry(result.ToEntry),
	}

	return res, nil
}

func validDataTransferRequest(req *pb.CreateTransferRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidDataID(req.GetFromAccountId()); err != nil {
		violations = append(violations, fieldViolation("from_account_id", err))
	}

	if err := val.ValidDataID(req.GetToAccountId()); err != nil {
		violations = append(violations, fieldViolation("to_account_id", err))
	}

	if err := val.ValidDateCurrency(req.GetCurrency()); err != nil {
		violations = append(violations, fieldViolation("currency", err))
	}

	if err := val.ValidDataAmount(req.GetAmount()); err != nil {
		violations = append(violations, fieldViolation("amount", err))
	}

	return violations
}

func (s *Server) validAccount(ctx context.Context, accountID int64, currency string) (db.Account, error) {
	account, err := s.store.GetAccount(ctx, accountID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return account, status.Errorf(codes.NotFound, "not found account %s", err)
		} else {
			return account, status.Errorf(codes.Internal, "failed to get account: %s", err)
		}
	}

	if account.Currency != currency {
		return account, status.Errorf(codes.InvalidArgument, "account (%d) currency not match: %s vs %s", account.ID, account.Currency, currency)
	}

	return account, nil
}
