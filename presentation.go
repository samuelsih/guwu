package guwu

import (
	"context"
	"errors"
	"net/http"

	"github.com/bytedance/sonic/decoder"
	"github.com/bytedance/sonic/encoder"
	"github.com/go-chi/chi/v5"
	"github.com/samuelsih/guwu/model"
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

func post[inType any, outType service.CommonOutput](sess model.SessionDeps, svc func(context.Context, *inType, model.Session) outType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userSession, err := lookupUser(sess, r)
		if err != nil {
			encoder.NewStreamEncoder(w).Encode(service.CommonResponse{
				StatusCode: http.StatusBadRequest,
				Msg:        err.Error(),
			})
			return
		}
		
		var in inType

		encoder := encoder.NewStreamEncoder(w)
		decoder := decoder.NewStreamDecoder(r.Body)
		decoder.DisallowUnknownFields()

		err = decoder.Decode(&in)
		if err != nil {
			encoder.Encode(service.CommonResponse{
				StatusCode: http.StatusBadRequest,
				Msg:        err.Error(),
			})

			return
		}

		out := svc(r.Context(), &in, *userSession)

		w.WriteHeader(out.Common().StatusCode)
		encoder.Encode(out)
	}
}

func put[inType any, outType service.CommonOutput](
	sess model.SessionDeps, 
	key string, 
	svc func(context.Context, string, *inType, model.Session) outType) http.HandlerFunc {
		
	return func(w http.ResponseWriter, r *http.Request) {
		userSess, err := lookupUser(sess, r)
		if err != nil {
			encoder.NewStreamEncoder(w).Encode(service.CommonResponse{
				StatusCode: http.StatusBadRequest,
				Msg:        err.Error(),
			})
			return
		}
		
		var in inType

		encoder := encoder.NewStreamEncoder(w)
		decoder := decoder.NewStreamDecoder(r.Body)
		decoder.DisallowUnknownFields()

		err = decoder.Decode(&in)
		if err != nil {
			encoder.Encode(service.CommonResponse{
				StatusCode: http.StatusBadRequest,
				Msg:        err.Error(),
			})

			return
		}

		out := svc(r.Context(), chi.URLParam(r, key), &in, *userSess)

		w.WriteHeader(out.Common().StatusCode)
		encoder.Encode(out)
	}
}

func delete[inType any, outType service.CommonOutput](sess model.SessionDeps, svc func(context.Context, *model.Session) outType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userSession, err := lookupUser(sess, r)
		if err != nil {
			encoder.NewStreamEncoder(w).Encode(service.CommonResponse{
				StatusCode: http.StatusBadRequest,
				Msg:        err.Error(),
			})
			return
		}

		out := svc(r.Context(), userSession)

		w.WriteHeader(out.Common().StatusCode)
		encoder.NewStreamEncoder(w).Encode(out)
	}
}

func loginOrRegister[inType any, outType service.CommonOutput](svc func(context.Context, *inType) outType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var in inType

		encoder := encoder.NewStreamEncoder(w)
		decoder := decoder.NewStreamDecoder(r.Body)
		decoder.DisallowUnknownFields()

		err := decoder.Decode(&in)
		if err != nil {
			encoder.Encode(service.CommonResponse{
				StatusCode: http.StatusBadRequest,
				Msg:        err.Error(),
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

func logout[outType service.CommonOutput](svc func(context.Context, string) outType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := r.Cookie("app_session")
		if err != nil {
			encoder.NewStreamEncoder(w).Encode(service.CommonResponse{
				StatusCode: http.StatusBadRequest,
				Msg:        err.Error(),
			})

			return
		}

		out := svc(r.Context(), sessionID.Value)

		w.WriteHeader(out.Common().StatusCode)
		encoder.NewStreamEncoder(w).Encode(out)
	}
}

func lookupUser(sessionDeps model.SessionDeps, r *http.Request) (*model.Session, error) {
	cookie, err := r.Cookie("app_session")
	if err != nil {
		return nil, errors.New(http.StatusText(http.StatusBadRequest))
	}

	user, err := sessionDeps.Get(r.Context(), cookie.Value)

	if err != nil || user == nil {
		return nil, errors.New(http.StatusText(http.StatusBadRequest))
	}

	return user, nil
}