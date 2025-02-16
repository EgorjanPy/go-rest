package save

import (
	"errors"
	"log/slog"
	resp "main/internal/lib/api/response"
	"main/internal/lib/logger/sl"
	"main/internal/lib/random"
	"main/internal/storage"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

const aliasLength = 6

type URLSaver interface {
	SaveUrl(urlToSave, alias string) (int64, error)
}
type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}
type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request")
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}
		log.Info("requesr body decoded", slog.Any("request", req))
		// err = validator.Validate.Struct(req)
		// validationErrors := err.(validator.ValidationErrors)
		// if err := validator.New().Struct(req); err != nil {
		// 	validateErr := err.(validator.ValidationErrors)
		// 	log.Error("invalid request", sl.Err(err))
		// 	render.JSON(w, r, resp.ValidationError(validateErr)) // Возвращает читабельную ошибку
		// 	return
		// }
		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)

		}
		id, err := urlSaver.SaveUrl(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("url", r.URL.Path))
			render.JSON(w, r, resp.Error("url already exists"))
			return
		}
		if err != nil {
			log.Error("failed to add url", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to add url"))
			return
		}
		log.Info("url added", slog.Int64("id", id))
		render.JSON(w, r, Response{
			Response: resp.OK(),
			Alias:    alias,
		})
	}
}
