package guwu

import (
	"context"
	"net/http"

	"github.com/bytedance/sonic/decoder"
	"github.com/bytedance/sonic/encoder"
	"github.com/go-chi/chi/v5"
	"github.com/samuelsih/guwu/service"
)

func POST[inType any, outType service.CommonOutput](mux *chi.Mux, route string, handler func(context.Context, *inType) outType) {
	mux.Post(route, func(w http.ResponseWriter, r *http.Request) {
		var in inType

		encoder := encoder.NewStreamEncoder(w)
		decoder := decoder.NewStreamDecoder(r.Body)
		decoder.DisallowUnknownFields()

		err := decoder.Decode(&in)
		if err != nil {
			encoder.Encode(service.CommonResponse{
				StatusCode: http.StatusBadRequest, 
				Msg: err.Error(),
			})

			return
		}

		out := handler(r.Context(), &in)

		if out.Common().StatusCode == 0 {
			w.WriteHeader(http.StatusOK)
			encoder.Encode(out)
			return
		}

		w.WriteHeader(out.Common().StatusCode)
		encoder.Encode(out)
	})
}

func PUT[inType any, outType service.CommonOutput](mux *chi.Mux, route string, param string, handler func(context.Context, *inType, string) outType) {
	mux.Post(route, func(w http.ResponseWriter, r *http.Request) {
		urlParam := chi.URLParam(r, param)
		
		var in inType

		encoder := encoder.NewStreamEncoder(w)
		decoder := decoder.NewStreamDecoder(r.Body)
		decoder.DisallowUnknownFields()

		err := decoder.Decode(&in)
		if err != nil {
			encoder.Encode(service.CommonResponse{
				StatusCode: http.StatusBadRequest, 
				Msg: err.Error(),
			})

			return
		}

		out := handler(r.Context(), &in, urlParam)

		if out.Common().StatusCode == 0 {
			w.WriteHeader(http.StatusOK)
			encoder.Encode(out)
			return
		}

		w.WriteHeader(out.Common().StatusCode)
		encoder.Encode(out)
	})
}

func DELETE[inType any, outType service.CommonOutput](mux *chi.Mux, route string, param string, handler func(context.Context, *inType, string) outType) {
	mux.Post(route, func(w http.ResponseWriter, r *http.Request) {
		urlParam := chi.URLParam(r, param)
		
		var in inType

		encoder := encoder.NewStreamEncoder(w)
		decoder := decoder.NewStreamDecoder(r.Body)
		decoder.DisallowUnknownFields()

		err := decoder.Decode(&in)
		if err != nil {
			encoder.Encode(service.CommonResponse{
				StatusCode: http.StatusBadRequest, 
				Msg: err.Error(),
			})

			return
		}

		out := handler(r.Context(), &in, urlParam)

		if out.Common().StatusCode == 0 {
			w.WriteHeader(http.StatusOK)
			encoder.Encode(out)
			return
		}

		w.WriteHeader(out.Common().StatusCode)
		encoder.Encode(out)
	})
}