package uploadOrder

import (
	"context"
	"errors"
	"net/http"

	"github.com/VladSnap/gopherLoyalty/internal/domain"
	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/api"
	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/log"
	"github.com/google/uuid"
	"github.com/swaggest/rest"
	"github.com/swaggest/usecase/status"
)

type UploadOrderUseCaseImpl struct {
	cmdHandler UploadOrderCmdHandler
}

func NewUploadOrderUseCase(cmdHandler UploadOrderCmdHandler) *UploadOrderUseCaseImpl {
	return &UploadOrderUseCaseImpl{cmdHandler: cmdHandler}
}

type UploadOrderCmdHandler interface {
	Execute(ctx context.Context, orderNumber string, currentUser uuid.UUID) (bool, error)
}

func (uc *UploadOrderUseCaseImpl) Execute(ctx context.Context, input UploadOrderRequest, output *api.HttpStatusResponse) error {
	currentUserID, ok := ctx.Value(api.KeyContext("UserID")).(uuid.UUID)
	if !ok {
		err := errors.New("current userID is empty")
		log.Zap.Error(err)
		return status.Wrap(err, status.Unknown)
	}

	_, err := uc.cmdHandler.Execute(ctx, input.TextBody, currentUserID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidOrderNumber):
			log.Zap.Error(err)
			return rest.HTTPCodeAsError(http.StatusUnprocessableEntity)
		case errors.Is(err, domain.ErrAlreadyUploadedOrderCurrent):
			log.Zap.Error(err)
			output.WithHttpStatus(http.StatusOK)
			return nil
		case errors.Is(err, domain.ErrAlreadyUploadedOrderAnother):
			log.Zap.Error(err)
			return status.Wrap(errors.New("order already uploaded by another user"), status.AlreadyExists)
		default:
			log.Zap.Error(err)
			return status.Wrap(errors.New("unknown error"), status.Unknown)
		}
	}
	output.WithHttpStatus(http.StatusAccepted)
	return nil
}
