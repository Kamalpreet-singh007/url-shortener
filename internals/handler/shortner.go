package handler
import(
	"net/http"
	"encoding/json"
	"net/url"
)

type shortenRequest struct{
	URL string `json:"url"`

}

type shortenResponse struct{
	ShortCode string `json:"short_code"`
}

func (h *UrlHandler) Shorten(w http.ResponseWriter, r *http.Request) {
	var req shortenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.Logger.Error("Failed to decode request body", "error", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return 
	}
	if req.URL == "" {
		h.Logger.Warn("Missing url field in payload")
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}
	_,err = url.ParseRequestURI(req.URL)
	if err != nil{	
		h.Logger.Warn("Invalid URL provided", "url", req.URL, "error", err)
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	
	URL, err :=h.Store.CreateUrl(r.Context(), req.URL)
	if err != nil{
		h.Logger.Error("Failed to create short URL", "error", err)
		http.Error(w, "Failed to create short URL", http.StatusInternalServerError)
		return
	}

	h.Logger.Info("Short URL created successfully", "short_code", URL.ShortCode, "original_url", URL.OriginalURL)
	res:= shortenResponse{ShortCode: URL.ShortCode}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}	
