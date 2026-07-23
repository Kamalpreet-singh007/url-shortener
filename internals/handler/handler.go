package handler

import(
	"github.com/Kamalpreet-singh007/url-shortener/internals/store"
	"log/slog"


)

type UrlHandler struct {
	Store store.URLStore
	Logger *slog.Logger
}

func NewUrlHandler(store store.URLStore, logger *slog.Logger) *UrlHandler {
	return &UrlHandler{
		Store: store,
		Logger: logger,
	}
}