package uploadOrder

import (
	"context"

	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/api"
)

type UploadOrderUseCaseImpl struct {
}

func NewUploadOrderUseCase() *UploadOrderUseCaseImpl {
	return &UploadOrderUseCaseImpl{}
}

func (uc *UploadOrderUseCaseImpl) Execute(ctx context.Context, input string, output *api.EmptyBody) error {
	output = &api.EmptyBody{}
	return nil
}
