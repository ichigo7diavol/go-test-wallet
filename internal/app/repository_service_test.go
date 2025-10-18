package app_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/ichigo7diabol/go-test-wallet/internal/app"
	"github.com/ichigo7diabol/go-test-wallet/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockWalletRepository — мок реализация интерфейса WalletRepositoryService
type MockWalletRepository struct {
	mock.Mock
}

func (m *MockWalletRepository) Create(initialBalance float32) (*models.WalletModel, error) {
	args := m.Called(initialBalance)
	if model, ok := args.Get(0).(*models.WalletModel); ok {
		return model, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockWalletRepository) GetByID(id uuid.UUID) (*models.WalletModel, error) {
	args := m.Called(id)
	if model, ok := args.Get(0).(*models.WalletModel); ok {
		return model, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockWalletRepository) UpdateBalance(id uuid.UUID, balance float32) (float32, float32, *models.WalletModel, error) {
	args := m.Called(id, balance)
	if model, ok := args.Get(2).(*models.WalletModel); ok {
		return args.Get(0).(float32), args.Get(1).(float32), model, args.Error(3)
	}
	return args.Get(0).(float32), args.Get(1).(float32), nil, args.Error(3)
}

func (m *MockWalletRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockWalletRepository) List() ([]models.WalletModel, error) {
	args := m.Called()
	if list, ok := args.Get(0).([]models.WalletModel); ok {
		return list, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockWalletRepository) Deposit(id uuid.UUID, amount float32) (float32, float32, *models.WalletModel, error) {
	args := m.Called(id, amount)
	if model, ok := args.Get(2).(*models.WalletModel); ok {
		return args.Get(0).(float32), args.Get(1).(float32), model, args.Error(3)
	}
	return args.Get(0).(float32), args.Get(1).(float32), nil, args.Error(3)
}

func (m *MockWalletRepository) Withdraw(id uuid.UUID, amount float32) (float32, float32, *models.WalletModel, error) {
	args := m.Called(id, amount)
	if model, ok := args.Get(2).(*models.WalletModel); ok {
		return args.Get(0).(float32), args.Get(1).(float32), model, args.Error(3)
	}
	return args.Get(0).(float32), args.Get(1).(float32), nil, args.Error(3)
}

func TestCreateWallet(t *testing.T) {
	repo := new(MockWalletRepository)
	service := app.NewWalletService(repo)

	expected := &models.WalletModel{Balance: 100}
	repo.On("Create", float32(100)).Return(expected, nil)

	wallet, err := service.CreateWallet(100)

	assert.NoError(t, err)
	assert.Equal(t, expected, wallet)
	repo.AssertExpectations(t)
}

func TestGetWallet(t *testing.T) {
	repo := new(MockWalletRepository)
	service := app.NewWalletService(repo)
	id := uuid.New()

	expected := &models.WalletModel{ID: id, Balance: 50}
	repo.On("GetByID", id).Return(expected, nil)

	wallet, err := service.GetWallet(id)

	assert.NoError(t, err)
	assert.Equal(t, expected, wallet)
	repo.AssertExpectations(t)
}

func TestDeleteWallet(t *testing.T) {
	repo := new(MockWalletRepository)
	service := app.NewWalletService(repo)
	id := uuid.New()

	repo.On("Delete", id).Return(nil)

	err := service.DeleteWallet(id)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestListWallets(t *testing.T) {
	repo := new(MockWalletRepository)
	service := app.NewWalletService(repo)

	list := []models.WalletModel{
		{Balance: 10},
		{Balance: 20},
	}
	repo.On("List").Return(list, nil)

	result, err := service.ListWallets()

	assert.NoError(t, err)
	assert.Equal(t, list, result)
	repo.AssertExpectations(t)
}

func TestChangeBalance_Deposit(t *testing.T) {
	repo := new(MockWalletRepository)
	service := app.NewWalletService(repo)
	id := uuid.New()

	old := float32(50)
	new := float32(100)
	wallet := &models.WalletModel{ID: id, Balance: new}

	repo.On("Withdraw", id, float32(50)).Return(old, new, wallet, nil)

	oldBalance, newBalance, model, err := service.ChangeBalance(id, app.DepositOperation, 50)

	assert.NoError(t, err)
	assert.Equal(t, old, oldBalance)
	assert.Equal(t, new, newBalance)
	assert.Equal(t, wallet, model)
	repo.AssertExpectations(t)
}

func TestChangeBalance_Withdraw(t *testing.T) {
	repo := new(MockWalletRepository)
	service := app.NewWalletService(repo)
	id := uuid.New()

	old := float32(100)
	new := float32(50)
	wallet := &models.WalletModel{ID: id, Balance: new}

	repo.On("Deposit", id, float32(50)).Return(old, new, wallet, nil)

	oldBalance, newBalance, model, err := service.ChangeBalance(id, app.WithdrawOperation, 50)

	assert.NoError(t, err)
	assert.Equal(t, old, oldBalance)
	assert.Equal(t, new, newBalance)
	assert.Equal(t, wallet, model)
	repo.AssertExpectations(t)
}

func TestChangeBalance_UnknownOperation(t *testing.T) {
	repo := new(MockWalletRepository)
	service := app.NewWalletService(repo)
	id := uuid.New()

	_, _, _, err := service.ChangeBalance(id, "INVALID_OP", 100)

	assert.ErrorIs(t, err, app.ErrUnknownOperation)
}
