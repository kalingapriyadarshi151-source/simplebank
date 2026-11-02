package db

import (
	"context"
	// "database/sql"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/techschool/simplebank/util"
)


func TestCreateEntry(t *testing.T) {
	CreateRandomAccount(t)
}

func TestGetEntry(t *testing.T) {

	account := CreateRandomAccount(t)

	args := CreateEntryParams{
		AccountID: account.ID,
		Amount: util.RandInt(0,2334342),
	}

	entry1 ,err := testQueries.CreateEntry(context.Background(), args) 
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)

	

	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry2.AccountID, args.AccountID)
	require.Equal(t, entry2.Amount, args.Amount)
	require.Equal(t, entry2.ID, entry1.ID)

}