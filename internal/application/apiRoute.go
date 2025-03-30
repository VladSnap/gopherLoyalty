package application

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/swaggest/rest/nethttp"
	"github.com/swaggest/usecase"
)

func regRoute[Input, Output any](router chi.Router,
	useCase func(ctx context.Context, input Input, output *Output) error,
	method string,
	path string,
	successStatus int,
	summary string,
	description string,
	tags ...string) {
	// Создание Interactor
	inter := usecase.NewInteractor(useCase)
	inter.SetTitle(summary)
	inter.SetDescription(description)
	inter.SetTags(tags...)

	var handler *nethttp.Handler

	if successStatus > 0 {
		handler = nethttp.NewHandler(inter, nethttp.SuccessStatus(successStatus))
	} else {
		handler = nethttp.NewHandler(inter)
	}

	router.Method(method, path, handler)
}
