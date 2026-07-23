package handler

import(
	"net/http"
)


func (h *UrlHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	shortcode :=r.PathValue("code")

	url, err := h.Store.GetByShortCode(r.Context(), shortcode)
	if errors.Is(err, sql.ErrNoRows) {
		h.Logger.Info("Short_code not found", "short_code", shortcode)
		http.Error(w, "Short_code not found", http.StatusNotFound)
		return
	}
	if err != nil {
		h.Logger.Info("Failed to look up  short_code", "short_code", shortcode, "error", err)
		http.Error(w, "Failed to look up short_code", http.InternalServerError)
		return
	}
	
	h.Logger.Info("Redirected successfully", "orignal_URL", url.OriginalURL)
	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}	