package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/makuo12/ghost_server/tools"
	"github.com/makuo12/ghost_server/utils"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	config, err := utils.LoadConfig("../..")
	require.NoError(t, err)
	user := createRandomUser(t)
	// We would test as money comes in Naira
	money := utils.RandomMoney()
	fmt.Println("money: ", money)
	// We create a user account for the user
	// We create dollar account
	moneyUSD, err := tools.ConvertPrice(money, utils.NGN, utils.USD, "1000", "1.38", uuid.New())
	fmt.Println("moneyUSD: ", moneyUSD)
	require.NoError(t, err)
	require.NotEmpty(t, moneyUSD)
	require.Equal(t, moneyUSD, utils.TestNairaToDollarFloat())

	accountUSD, err := testStore.CreateTestAccount(context.Background(), CreateTestAccountParams{
		UserID:   user.ID,
		Balance:  tools.MoneyFloatToInt(moneyUSD),
		Currency: utils.USD,
	})
	require.NoError(t, err)
	require.NotEmpty(t, accountUSD)
	require.Equal(t, accountUSD.Balance, utils.TestNairaToDollarInt())
	require.Equal(t, accountUSD.UserID, user.ID)
	require.Equal(t, accountUSD.Currency, utils.USD)
	// We create naira account
	moneyNGN, err := tools.ConvertPrice(money, utils.NGN, utils.NGN, config.DollarToNaira, config.DollarToCAD, uuid.New())
	require.NoError(t, err)
	require.NotEmpty(t, moneyNGN)
	require.Equal(t, moneyNGN, utils.TestNairaToNairaFloat())
	accountNGN, err := testStore.CreateTestAccount(context.Background(), CreateTestAccountParams{
		UserID:   user.ID,
		Balance:  tools.MoneyFloatToInt(moneyNGN),
		Currency: utils.NGN,
	})
	require.NoError(t, err)
	require.NotEmpty(t, accountNGN)
	require.Equal(t, accountNGN.Balance, utils.TestNairaToNairaInt())
	require.Equal(t, accountNGN.UserID, user.ID)
	require.Equal(t, accountNGN.Currency, utils.NGN)

	return accountNGN
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}
