package app

import (
	"errors"

	"github.com/google/uuid"
	"github.com/ichigo7diabol/go-test-wallet/internal/models"
)

type WalletOperation string

const (
	WithdrawOperation WalletOperation = "WITHDRAW"
	DepositOperation  WalletOperation = "DEPOSIT"
)

var (
	ErrUnknownOperation = errors.New("unknown operation")
)

type WalletService struct {
	repository WalletRepositoryService
}

func NewWalletService(repo WalletRepositoryService) *WalletService {
	return &WalletService{repository: repo}
}

func (s *WalletService) CreateWallet(initialBalance float32) (*models.WalletModel, error) {
	return s.repository.Create(initialBalance)
}

func (s *WalletService) GetWallet(id uuid.UUID) (*models.WalletModel, error) {
	return s.repository.GetByID(id)
}

func (s *WalletService) DeleteWallet(id uuid.UUID) error {
	return s.repository.Delete(id)
}

func (s *WalletService) ListWallets() ([]models.WalletModel, error) {
	return s.repository.List()
}

func (s *WalletService) ChangeBalance(id uuid.UUID, op WalletOperation, amount float32) (
	oldBalance float32,
	newBalance float32,
	model *models.WalletModel,
	err error,
) {
	switch op {
	case WithdrawOperation:
		return s.repository.Deposit(id, amount)
	case DepositOperation:
		return s.repository.Withdraw(id, amount)
	default:
		return 0, 0, nil, ErrUnknownOperation
	}
}
