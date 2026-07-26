package handler

import(
	"github.com/Kamalpreet-singh007/url-shortener/internals/store"


)

type UrlHandler struct {
	Store store.URLStore
}

func NewUrlHandler(store store.URLStore) *UrlHandler {
	return &UrlHandler{
		Store: store,
	}
}