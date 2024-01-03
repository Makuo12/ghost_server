package db

import (
	"context"
	"flex_server/tools"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store provides all functions to execute db queries and transactions

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	connPool *pgxpool.Pool
	*Queries
}

// NewStore creates a new store

func NewStore(connPool *pgxpool.Pool) *SQLStore {
	return &SQLStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}

// execTx executes a function within a database transaction

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID    uuid.UUID `json:"from_account_id"`
	ToAccountID      uuid.UUID `json:"to_account_id"`
	FromAccountIDInt int64     `json:"from_account_id_int"`
	ToAccountIDInt   int64     `json:"to_account_id_int"`
	Amount           string    `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var txKey = struct{}{}

// TransferTx performs a money transfer from one account to the other.
// It creates a transfer record, add account entries, and update accounts' balance within a single database transaction
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		//txName := ctx.Value(txKey)
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID:    arg.FromAccountID,
			ToAccountID:      arg.ToAccountID,
			FromAccountIDInt: arg.FromAccountIDInt,
			ToAccountIDInt:   arg.ToAccountIDInt,
			Amount:           tools.MoneyStringToInt(arg.Amount),
		})
		if err != nil {
			return err
		}
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -tools.MoneyStringToInt(arg.Amount),
		})
		if err != nil {
			return err
		}
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    tools.MoneyStringToInt(arg.Amount),
		})
		if err != nil {
			return err
		}
		// TODO update accounts' balance

		if arg.FromAccountIDInt < arg.ToAccountIDInt {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, "-"+arg.Amount, arg.ToAccountID, arg.Amount)
			if err != nil {
				return err
			}
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, "-"+arg.Amount)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return result, err

}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountIDOne uuid.UUID,
	amountOne string,
	accountIDTwo uuid.UUID,
	amountTwo string,
) (accountOne Account, accountTwo Account, err error) {
	accountOne, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountIDOne,
		Amount: tools.MoneyStringToInt(amountOne),
	})
	if err != nil {
		return
	}
	accountTwo, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountIDTwo,
		Amount: tools.MoneyStringToInt(amountTwo),
	})
	if err != nil {
		return
	}
	return
}
