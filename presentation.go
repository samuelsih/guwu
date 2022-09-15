package guwu

import (
	"context"
	"net/http"

	"github.com/bytedance/sonic/decoder"
	"github.com/bytedance/sonic/encoder"
	"github.com/go-chi/chi/v5"
	"github.com/samuelsih/guwu/service"
)

func get[outType service.CommonOutput](svc func(context.Context) outType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		out := svc(r.Context())

		if out.Common().StatusCode == 0 {
			encoder := encoder.NewStreamEncoder(w)
			
			w.WriteHeader(http.StatusOK)
			encoder.Encode(out)
			return
		}

		encoder := encoder.NewStreamEncoder(w)

		w.WriteHeader(out.Common().StatusCode)
		encoder.Encode(out)
	}
}

func getWithParam[outType service.CommonOutput](svc func(context.Context, string) outType, param string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlParam := chi.URLParam(r, param)	

		out := svc(r.Context(), urlParam)

		if out.Common().StatusCode == 0 {
			encoder := encoder.NewStreamEncoder(w)
			
			w.WriteHeader(http.StatusOK)
			encoder.Encode(out)
			return
		}

		encoder := encoder.NewStreamEncoder(w)

		w.WriteHeader(out.Common().StatusCode)
		encoder.Encode(out)
	}
}

func post[inType any, outType service.CommonOutput](svc func(context.Context, *inType) outType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		out := svc(r.Context(), &in)

		if out.Common().StatusCode == 0 {
			w.WriteHeader(http.StatusOK)
			encoder.Encode(out)
			return
		}

		w.WriteHeader(out.Common().StatusCode)
		encoder.Encode(out)
	}
}

func put[inType any, outType service.CommonOutput](svc func(context.Context, *inType, string) outType, param string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		urlParam := chi.URLParam(r, param)

		out := svc(r.Context(), &in, urlParam)

		if out.Common().StatusCode == 0 {
			w.WriteHeader(http.StatusOK)
			encoder.Encode(out)
			return
		}

		w.WriteHeader(out.Common().StatusCode)
		encoder.Encode(out)
	}
}