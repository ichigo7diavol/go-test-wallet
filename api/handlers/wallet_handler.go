package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/ichigo7diabol/go-test-wallet/api/openapi"
	"github.com/ichigo7diabol/go-test-wallet/internal/app"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

var (
	ErrIncorrectData  = errors.New("incorrect data")
	ErrInternalServer = errors.New("internal server error")
)

type WalletHandler struct {
	WalletService *app.WalletService
}

func NewWalletHandler(walletService *app.WalletService) *WalletHandler {
	return &WalletHandler{
		WalletService: walletService,
	}
}

func (h *WalletHandler) CreateWallet(ctx echo.Context) error {
	var req openapi.CreateWalletRequest
	if err := ctx.Bind(&req); err != nil {
		return NewHttpError(err)
	}
	model, err := h.WalletService.CreateWallet(req.InitialBalance)
	if err != nil {
		return NewHttpError(err)
	}
	wallet := openapi.Wallet{
		WalletId:  (*openapi_types.UUID)(&model.ID),
		Balance:   &model.Balance,
		CreatedAt: &model.CreatedAt,
		UpdatedAt: &model.UpdatedAt,
	}
	return ctx.JSON(http.StatusCreated, wallet)
}

// Операции с кошельком (DEPOSIT/WITHDRAW)
func (h *WalletHandler) ChangeWallet(ctx echo.Context) error {
	var req openapi.WalletOperationRequest
	if err := ctx.Bind(&req); err != nil {
		return NewHttpError(err)
	}
	oldBalance, newBalance, model, err := h.WalletService.ChangeBalance(
		uuid.UUID(req.WalletId),
		app.WalletOperation(req.OperationType),
		req.Amount,
	)
	if err != nil {
		return NewHttpError(err)
	}
	now := time.Now()
	resp := openapi.WalletOperationResponse{
		WalletId:      &model.ID,
		OperationType: (*openapi.WalletOperationResponseOperationType)(&req.OperationType),
		OldBalance:    &oldBalance,
		NewBalance:    &newBalance,
		Amount:        &req.Amount,
		Timestamp:     &now,
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (h *WalletHandler) ListWallets(ctx echo.Context) error {
	models, err := h.WalletService.ListWallets()
	if err == nil {
		return NewHttpError(err)
	}
	wallets := make([]openapi.Wallet, len(models))
	for i, model := range models {
		wallets[i] = openapi.Wallet{
			WalletId:  (*openapi_types.UUID)(&model.ID),
			Balance:   &model.Balance,
			CreatedAt: &model.CreatedAt,
			UpdatedAt: &model.UpdatedAt,
		}
	}
	return ctx.JSON(http.StatusOK, wallets)
}

func (h *WalletHandler) GetWallet(ctx echo.Context, walletId openapi_types.UUID) error {
	model, err := h.WalletService.GetWallet(walletId)
	if err != nil {
		return NewHttpError(err)
	}
	return ctx.JSON(http.StatusOK, openapi.Wallet{
		WalletId:  &walletId,
		Balance:   &model.Balance,
		CreatedAt: &model.CreatedAt,
		UpdatedAt: &model.UpdatedAt,
	})
}

func (h *WalletHandler) DeleteWallet(ctx echo.Context, walletId openapi_types.UUID) error {
	err := h.WalletService.DeleteWallet(walletId)
	if err != nil {
		return NewHttpError(err)
	}
	return ctx.NoContent(http.StatusNoContent)
}

func NewHttpError(err error) error {
	var code = http.StatusInternalServerError
	var msg = map[string]string{}

	switch {
	case errors.Is(err, app.ErrInvalidAmount):
	case errors.Is(err, app.ErrUnknownOperation):
	case errors.Is(err, ErrIncorrectData):
		code = http.StatusBadRequest
	case errors.Is(err, app.ErrInsufficientFunds):
		code = http.StatusPaymentRequired
	case errors.Is(err, app.ErrWalletNotFound):
		code = http.StatusNotFound
	default:
		err = ErrInternalServer
		code = http.StatusInternalServerError
	}
	msg = map[string]string{"error": err.Error()}

	return echo.NewHTTPError(code, msg)
}
