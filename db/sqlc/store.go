package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

// creates a new store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// transferTxParams contais the input parameters of transfer Transaction
type transferTxParams struct {
	FromAccountId int64 `json:"from_accout_id"`
	ToAccountId   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// transferTxResults contains the result of transfer Transaction
type transferTxResults struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var txKey = struct{}{}

// executes a func in a single database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {

	tx, err := store.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	q := New(tx)

	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("err: %v and the rbErr: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// Transfertx performs a money transfer from one account to another
// it creates a transfer record , add account entries and update account's balance within a single database
func (store *Store) Transfertx(ctx context.Context, arg transferTxParams) (transferTxResults, error) {
	var result transferTxResults

	err := store.execTx(ctx, func(q *Queries) error {

		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountId,
			ToAccountID:   arg.ToAccountId,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		
	
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountId,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

	
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountId,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// todo: update account's balance

		if arg.FromAccountId < arg.ToAccountId {
			result.FromAccount , result.ToAccount, err = AddMoney(ctx, q , arg.FromAccountId, -arg.Amount, arg.ToAccountId, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount ,err = AddMoney(ctx, q , arg.ToAccountId, arg.Amount, arg.FromAccountId, -arg.Amount)
		}

		return nil
	})

	return result, err
}

func AddMoney(
	ctx context.Context,
	q *Queries,
	account1ID int64,
	amount1 int64,
	account2ID int64,
	amount2 int64,
) (account1 Account, account2 Account , err error) {

	account1, err = q.AddAccountBalance(ctx , AddAccountBalanceParams{
		ID: account1ID,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx , AddAccountBalanceParams{
		ID: account2ID,
		Amount: amount2,
	})

	return
}