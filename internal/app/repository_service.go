package app

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ichigo7diabol/go-test-wallet/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrInvalidAmount     = errors.New("invalid amount")
	ErrWalletNotFound    = errors.New("wallet not found")
)

type RepositoryService struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *RepositoryService {
	return &RepositoryService{db: db}
}

type WalletRepositoryService interface {
	Create(initialBalance float32) (*models.WalletModel, error)
	GetByID(id uuid.UUID) (*models.WalletModel, error)
	UpdateBalance(id uuid.UUID, balance float32) (oldBalance float32, newBalance float32, model *models.WalletModel, err error)
	Delete(id uuid.UUID) error
	List() ([]models.WalletModel, error)
	Deposit(id uuid.UUID, amount float32) (oldBalance float32, newBalance float32, model *models.WalletModel, err error)
	Withdraw(id uuid.UUID, amount float32) (oldBalance float32, newBalance float32, model *models.WalletModel, err error)
}

func (r *RepositoryService) Create(initialBalance float32) (*models.WalletModel, error) {
	w := &models.WalletModel{
		Balance:   initialBalance,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := r.db.Create(w).Error; err != nil {
		return nil, err
	}
	return w, nil
}

func (r *RepositoryService) GetByID(id uuid.UUID) (*models.WalletModel, error) {
	var w models.WalletModel
	if err := r.db.First(&w, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrWalletNotFound
		}
		return nil, err
	}
	return &w, nil
}

func (r *RepositoryService) UpdateBalance(id uuid.UUID, balance float32) (
	oldBalance float32,
	newBalance float32,
	model *models.WalletModel,
	err error,
) {
	if balance < 0 {
		return 0, 0, nil, ErrInvalidAmount
	}
	var w models.WalletModel
	err = r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&w, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrWalletNotFound
			}
			return err
		}
		oldBalance = w.Balance

		w.Balance = balance
		w.UpdatedAt = time.Now()

		newBalance = w.Balance

		if err := tx.Save(&w).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return 0, 0, nil, err
	}
	return oldBalance, newBalance, &w, nil
}

func (r *RepositoryService) Delete(id uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var w models.WalletModel
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&w, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrWalletNotFound
			}
			return err
		}

		return tx.Delete(&w).Error
	})
}

func (r *RepositoryService) List() ([]models.WalletModel, error) {
	var wallets []models.WalletModel
	if err := r.db.Find(&wallets).Error; err != nil {
		return nil, err
	}
	return wallets, nil
}

func (r *RepositoryService) Deposit(id uuid.UUID, amount float32) (
	oldBalance float32,
	newBalance float32,
	model *models.WalletModel,
	err error,
) {
	if amount < 0 {
		return 0, 0, nil, ErrInvalidAmount
	}
	var w models.WalletModel

	err = r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&w, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrWalletNotFound
			}
			return err
		}
		oldBalance = w.Balance

		w.Balance += amount
		w.UpdatedAt = time.Now()

		newBalance = w.Balance

		return tx.Save(&w).Error
	})

	if err != nil {
		return 0, 0, nil, err
	}
	return oldBalance, newBalance, &w, nil
}

func (r *RepositoryService) Withdraw(id uuid.UUID, amount float32) (
	oldBalance float32,
	newBalance float32,
	model *models.WalletModel,
	err error,
) {
	if amount < 0 {
		return 0, 0, nil, ErrInvalidAmount
	}
	var w models.WalletModel

	err = r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&w, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrWalletNotFound
			}
			return err
		}

		if w.Balance < amount {
			return ErrInsufficientFunds
		}

		oldBalance = w.Balance

		w.Balance -= amount
		w.UpdatedAt = time.Now()

		newBalance = w.Balance

		return tx.Save(&w).Error
	})

	if err != nil {
		return 0, 0, nil, err
	}
	return oldBalance, newBalance, &w, nil
}
