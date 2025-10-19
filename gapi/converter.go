package gapi

import (
	db "github.com/JeongWoo-Seo/simpleBank/db/sqlc"
	"github.com/JeongWoo-Seo/simpleBank/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func converterUser(user db.User) *pb.User {
	return &pb.User{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
	}
}

func converterAccount(account db.Account) *pb.Account {
	return &pb.Account{
		Id:        account.ID,
		Owner:     account.Owner,
		Balance:   account.Balance,
		Currency:  account.Currency,
		CreatedAt: timestamppb.New(account.CreatedAt),
	}
}

func converterTransfer(transfer db.Transfer) *pb.Transfer {
	return &pb.Transfer{
		Id:            transfer.ID,
		FromAccountId: transfer.FromAccountID,
		ToAccountId:   transfer.ToAccountID,
		Amount:        transfer.Amount,
		CreatedAt:     timestamppb.New(transfer.CreatedAt),
	}
}

func converterEntry(entry db.Entry) *pb.Entry {
	return &pb.Entry{
		Id:        entry.ID,
		AccountId: entry.AccountID,
		Amount:    entry.Amount,
		CreatedAt: timestamppb.New(entry.CreatedAt),
	}
}
