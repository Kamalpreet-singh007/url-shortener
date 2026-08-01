package handler

import(
	"net/http"
	"errors"
	"github.com/Kamalpreet-singh007/url-shortener/internals/store"
)


func (h *UrlHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	shortcode :=r.PathValue("code")

	url, err := h.Store.GetByShortCode(r.Context(), shortcode)
	if errors.Is(err, store.ErrNotFound) {
		http.Error(w, "Short_code not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed to look up short_code", http.StatusInternalServerError)
		return
	}
	
	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}	