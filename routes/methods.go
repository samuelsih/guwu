//UNUSED. Maybe will be used in the future
package routes

import (
	"context"
	"encoding/json"

	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/samuelsih/guwu/service"
)

type (
	PostHandler [T service.CommonInput, U service.CommonOutput] func(context.Context, *T) U
	GetHandler [U service.CommonOutput] func(context.Context) U
	GetWithParamHandler [U service.CommonOutput] func(context.Context, string) U
	GetWithQueryParamHandler [T service.CommonInput, U service.CommonOutput] func(context.Context, T, map[string][]string) U
) 

func Get[out service.CommonOutput](path string, handler GetHandler[out]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		output := handler(r.Context())
		
		encoder := json.NewEncoder(w)

		if output.CommonRes().StatusCode == 0 {
			w.WriteHeader(http.StatusOK)
			encoder.Encode(output)
			return
		}

		w.WriteHeader(output.CommonRes().StatusCode)
		encoder.Encode(output)
	}
}

func GetWithParam[out service.CommonOutput](path, param string, handler GetWithParamHandler[out]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlParam := chi.URLParam(r, param)
		output := handler(r.Context(), urlParam)
		encoder := json.NewEncoder(w)

		if output.CommonRes().StatusCode == 0 {
			w.WriteHeader(http.StatusOK)
			encoder.Encode(output)
			return
		}

		w.WriteHeader(output.CommonRes().StatusCode)
		encoder.Encode(output)
	}
}

func GetWithQueryParam[T service.CommonInput, U service.CommonOutput](path string, handler GetWithQueryParamHandler[T, U]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var in T
		var err error

		qp := r.URL.Query()
		
		in.CommonReq().UserSession, err = deps.readCookie(r)
		encoder := json.NewEncoder(w)

		defer r.Body.Close()

		if err != nil {
			encoder.Encode(service.CommonResponse{
				StatusCode: http.StatusBadRequest,
				Msg:        err.Error(),
			})
			return
		}

		output := handler(r.Context(), in, qp)

		if output.CommonRes().StatusCode == 0 {
			w.WriteHeader(http.StatusOK)
			encoder.Encode(output)
			return
		}

		w.WriteHeader(output.CommonRes().StatusCode)
		encoder.Encode(output)
	}
}

func Post[T service.CommonInput, U service.CommonOutput](path string, handler PostHandler[T, U]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userSession, err := deps.readCookie(r)

		encoder := json.NewEncoder(w)
		decoder := json.NewDecoder(r.Body)

		defer r.Body.Close()

		if err != nil {
			encoder.Encode(service.CommonResponse{
				StatusCode: http.StatusBadRequest,
				Msg:        err.Error(),
			})
			return
		}
		
		var in T

		in.CommonReq().UserSession = userSession

		decoder.DisallowUnknownFields()

		err = decoder.Decode(&in)
		if err != nil {
			encoder.Encode(service.CommonResponse{
				StatusCode: http.StatusBadRequest,
				Msg:        err.Error(),
			})

			return
		}

		out := handler(r.Context(), &in)

		if out.CommonRes().SessionID != "" {
			if errCookie := deps.setCookie(w, out.CommonRes().SessionID); errCookie != nil {
				w.WriteHeader(http.StatusInternalServerError)
				encoder.Encode(out)
				return
			}
		}

		w.WriteHeader(out.CommonRes().StatusCode)
		encoder.Encode(out)
	}
}