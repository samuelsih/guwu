package guwu

import (
	"context"
	"errors"
	"net/http"
	"os"
	"sync"

	"github.com/bytedance/sonic/decoder"
	"github.com/bytedance/sonic/encoder"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/securecookie"
	"github.com/samuelsih/guwu/model"
	"github.com/samuelsih/guwu/service"
)

var (
	securer = securecookie.New([]byte(os.Getenv("HASH_KEY")), nil)
	mutex   sync.RWMutex
)

func get[outType service.CommonOutput](svc func(context.Context) outType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		out := svc(r.Context())

		if out.CommonRes().StatusCode == 0 {
			encoder := encoder.NewStreamEncoder(w)

			w.WriteHeader(http.StatusOK)
			encoder.Encode(out)
			return
		}

		encoder := encoder.NewStreamEncoder(w)

		w.WriteHeader(out.CommonRes().StatusCode)
		encoder.Encode(out)
	}
}

func getWithParam[outType service.CommonOutput](svc func(context.Context, string) outType, param string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlParam := chi.URLParam(r, param)

		out := svc(r.Context(), urlParam)

		if out.CommonRes().StatusCode == 0 {
			encoder := encoder.NewStreamEncoder(w)

			w.WriteHeader(http.StatusOK)
			encoder.Encode(out)
			return
		}

		encoder := encoder.NewStreamEncoder(w)

		w.WriteHeader(out.CommonRes().StatusCode)
		encoder.Encode(out)
	}
}

func post[inType service.CommonInput, outType service.CommonOutput](
	sess model.SessionDeps, 
	svc func(context.Context, *inType) outType,
	) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userSession, err := readCookie(sess, r)
		if err != nil {
			encoder.NewStreamEncoder(w).Encode(service.CommonResponse{
				StatusCode: http.StatusBadRequest,
				Msg:        err.Error(),
			})
			return
		}
		
		var in inType

		in.CommonReq().UserSession = userSession

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

		out := svc(r.Context(), &in)

		if out.CommonRes().SessionID != "" {
			if errCookie := setCookie(w, out.CommonRes().SessionID); errCookie != nil {
				w.WriteHeader(http.StatusInternalServerError)
				encoder.Encode(out)
				return
			}
		}

		w.WriteHeader(out.CommonRes().StatusCode)
		encoder.Encode(out)
	}
}

func put[inType service.CommonInput, outType service.CommonOutput](
	sess model.SessionDeps, 
	key string, 
	svc func(context.Context, string, *inType) outType,
	) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userSession, err := readCookie(sess, r)
		if err != nil {
			encoder.NewStreamEncoder(w).Encode(service.CommonResponse{
				StatusCode: http.StatusBadRequest,
				Msg:        err.Error(),
			})
			return
		}
		
		var in inType

		in.CommonReq().UserSession = userSession

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

		out := svc(r.Context(), chi.URLParam(r, key), &in)

		w.WriteHeader(out.CommonRes().StatusCode)
		encoder.Encode(out)
	}
}

func delete[inType service.CommonInput, outType service.CommonOutput](
	sess model.SessionDeps, 
	key string, 
	svc func(context.Context, string) outType,
	) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {		
		userSession, err := readCookie(sess, r)
		if err != nil {
			encoder.NewStreamEncoder(w).Encode(service.CommonResponse{
				StatusCode: http.StatusBadRequest,
				Msg:        err.Error(),
			})
			return
		}
		
		var in inType

		in.CommonReq().UserSession = userSession

		out := svc(r.Context(), chi.URLParam(r, key))

		w.WriteHeader(out.CommonRes().StatusCode)
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

		if out.CommonRes().StatusCode >= 400 {
			w.WriteHeader(out.CommonRes().StatusCode)
			encoder.Encode(out)
			return
		}

		err = setCookie(w, out.CommonRes().SessionID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			encoder.Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(out.CommonRes().StatusCode)
		encoder.Encode(out)
	}
}

func logout[outType service.CommonOutput](svc func(context.Context, string) outType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("app_session")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			encoder.NewStreamEncoder(w).Encode(service.CommonResponse{
				StatusCode: http.StatusBadRequest,
				Msg:        `cookie not found`,
			})
			return
		}

		var sessionID string

		err = securer.Decode("app_session", cookie.Value, &sessionID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			encoder.NewStreamEncoder(w).Encode(service.CommonResponse{
				StatusCode: http.StatusBadRequest,
				Msg:        `unknown user`,
			})
			return
		}

		out := svc(r.Context(), sessionID)

		w.WriteHeader(out.CommonRes().StatusCode)
		encoder.NewStreamEncoder(w).Encode(out)
	}
}

func readCookie(sessionDeps model.SessionDeps, r *http.Request) (model.Session, error) {
	cookie, err := r.Cookie("app_session")
	if err != nil || cookie.Value == "" {
		return model.Session{}, errors.New(http.StatusText(http.StatusBadRequest))
	}

	var sessionID string

	err = securer.Decode("app_session", cookie.Value, &sessionID)
	if err != nil {
		return model.Session{}, errors.New(http.StatusText(http.StatusBadRequest))
	}

	user, err := sessionDeps.Get(r.Context(), sessionID)

	if err != nil {
		return model.Session{}, errors.New(http.StatusText(http.StatusBadRequest))
	}

	return user, nil
}

func setCookie(w http.ResponseWriter, value string) error {
	mutex.Lock()
	defer mutex.Unlock()

	encoded, err := securer.Encode("app_session", value)
	if err != nil {
		return err
	}

	cookie := http.Cookie{
		Name:     "app_session",
		Value:    encoded,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		MaxAge:   24 * 3600,
	}

	http.SetCookie(w, &cookie)
	return nil
}

