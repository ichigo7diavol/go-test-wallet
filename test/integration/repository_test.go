//go:build integration

package app_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ichigo7diabol/go-test-wallet/internal/app"
	"github.com/ichigo7diabol/go-test-wallet/internal/models"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.WalletModel{})
	require.NoError(t, err)

	return db
}

func TestRepository_CreateAndGet(t *testing.T) {
	db := setupTestDB(t)
	repo := app.NewRepository(db)

	w, err := repo.Create(100)
	require.NoError(t, err)
	require.NotNil(t, w)
	require.Equal(t, float32(100), w.Balance)

	got, err := repo.GetByID(w.ID)
	require.NoError(t, err)
	require.Equal(t, w.ID, got.ID)
	require.Equal(t, w.Balance, got.Balance)
}

func TestRepository_Create_InvalidAmount(t *testing.T) {
	db := setupTestDB(t)
	repo := app.NewRepository(db)

	w, err := repo.Create(-10)
	require.ErrorIs(t, err, app.ErrInvalidAmount)
	require.Nil(t, w)
}

func TestRepository_UpdateBalance(t *testing.T) {
	db := setupTestDB(t)
	repo := app.NewRepository(db)

	w, _ := repo.Create(50)
	old, newBal, updated, err := repo.UpdateBalance(w.ID, 200)
	require.NoError(t, err)
	require.Equal(t, float32(50), old)
	require.Equal(t, float32(200), newBal)
	require.Equal(t, float32(200), updated.Balance)
}

func TestRepository_UpdateBalance_Invalid(t *testing.T) {
	db := setupTestDB(t)
	repo := app.NewRepository(db)

	w, _ := repo.Create(50)
	_, _, _, err := repo.UpdateBalance(w.ID, -5)
	require.ErrorIs(t, err, app.ErrInvalidAmount)
}

func TestRepository_UpdateBalance_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := app.NewRepository(db)

	id := uuid.New()
	_, _, _, err := repo.UpdateBalance(id, 100)
	require.ErrorIs(t, err, app.ErrWalletNotFound)
}

func TestRepository_Deposit(t *testing.T) {
	db := setupTestDB(t)
	repo := app.NewRepository(db)

	w, _ := repo.Create(100)
	old, newBal, updated, err := repo.Deposit(w.ID, 50)
	require.NoError(t, err)
	require.Equal(t, float32(100), old)
	require.Equal(t, float32(150), newBal)
	require.Equal(t, float32(150), updated.Balance)
}

func TestRepository_Deposit_InvalidAmount(t *testing.T) {
	db := setupTestDB(t)
	repo := app.NewRepository(db)

	w, _ := repo.Create(100)
	_, _, _, err := repo.Deposit(w.ID, -5)
	require.ErrorIs(t, err, app.ErrInvalidAmount)
}

func TestRepository_Withdraw(t *testing.T) {
	db := setupTestDB(t)
	repo := app.NewRepository(db)

	w, _ := repo.Create(100)
	old, newBal, updated, err := repo.Withdraw(w.ID, 60)
	require.NoError(t, err)
	require.Equal(t, float32(100), old)
	require.Equal(t, float32(40), newBal)
	require.Equal(t, float32(40), updated.Balance)
}

func TestRepository_Withdraw_InsufficientFunds(t *testing.T) {
	db := setupTestDB(t)
	repo := app.NewRepository(db)

	w, _ := repo.Create(30)
	_, _, _, err := repo.Withdraw(w.ID, 50)
	require.ErrorIs(t, err, app.ErrInsufficientFunds)
}

func TestRepository_Withdraw_InvalidAmount(t *testing.T) {
	db := setupTestDB(t)
	repo := app.NewRepository(db)

	w, _ := repo.Create(30)
	_, _, _, err := repo.Withdraw(w.ID, -10)
	require.ErrorIs(t, err, app.ErrInvalidAmount)
}

func TestRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := app.NewRepository(db)

	w, _ := repo.Create(70)
	err := repo.Delete(w.ID)
	require.NoError(t, err)

	_, err = repo.GetByID(w.ID)
	require.ErrorIs(t, err, app.ErrWalletNotFound)
}

func TestRepository_Delete_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := app.NewRepository(db)

	id := uuid.New()
	err := repo.Delete(id)
	require.ErrorIs(t, err, app.ErrWalletNotFound)
}

func TestRepository_List(t *testing.T) {
	db := setupTestDB(t)
	repo := app.NewRepository(db)

	for i := 0; i < 3; i++ {
		_, err := repo.Create(float32(10 * i))
		require.NoError(t, err)
		time.Sleep(5 * time.Millisecond)
	}

	list, err := repo.List()
	require.NoError(t, err)
	require.Len(t, list, 3)
	require.GreaterOrEqual(t, list[1].CreatedAt.UnixNano(), list[0].CreatedAt.UnixNano())
}
