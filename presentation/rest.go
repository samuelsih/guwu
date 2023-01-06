package presentation

import (
	"context"
	"log"
	"net/http"

	b "github.com/samuelsih/guwu/business"
	"github.com/samuelsih/guwu/pkg/request"
	"github.com/samuelsih/guwu/pkg/response"
)

type DefaultHandler[out b.CommonOutput] func (ctx context.Context, commonIn b.CommonInput) out

type InputHandler[inType any, out b.CommonOutput] func (ctx context.Context, in inType, commonIn b.CommonInput) out

type Opts struct {
	GetSessionCookie bool
	SetSessionCookie bool
	DecodeRequestBody bool
	
	URLParams []string
	QueryParams []string
}

var (
	DefaultOpts = Opts{}

	GetterSetterSessionOpts = Opts {
		GetSessionCookie: true,
		SetSessionCookie: true,
	}

	GetSessionWithDecodeOpts = Opts {
		GetSessionCookie: true,
		DecodeRequestBody: true,
	}

	SetSessionWithDecodeOpts = Opts {
		SetSessionCookie: true,
		DecodeRequestBody: true,
	}

	OnlyDecodeOpts = Opts {
		DecodeRequestBody: true,
	}
)

func Get[outType b.CommonOutput](handle DefaultHandler[outType], opts Opts) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var commonInput = b.CommonInput{
			URLParam: opts.URLParams,
			QueryParam: opts.QueryParams,
		}

		var err error

		if opts.GetSessionCookie {
			commonInput.SessionID, err = getSessionCookie(r)
			if err != nil {
				encodeErr := response.JSON(w, 400, map[string]any{
					"code": 400,
					"messsage": err.Error(),
				})
	
				if encodeErr != nil {
					log.Printf("presentation.Get: %v", err)
				}
	
				return
			}
		}

		out := handle(r.Context(), commonInput)

		if err := response.JSON(w, out.CommonRes().StatusCode, &out); err != nil {
			log.Printf("presentation.Get: %v", err)
			return
		}
	}
}

func Post[inType any, outType b.CommonOutput](handle InputHandler[inType, outType], opts Opts) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var in inType
		var err error

		var commonInput = b.CommonInput{
			URLParam: opts.URLParams,
			QueryParam: opts.QueryParams,
		}

		if opts.GetSessionCookie {
			commonInput.SessionID, err = getSessionCookie(r)
			if err != nil {
				encodeErr := response.JSON(w, 400, map[string]any{
					"code": 400,
					"message": err.Error(),
				})
	
				if encodeErr != nil {
					log.Printf("presentation.Get: %v", err)
				}
	
				return
			}
		}

		if opts.DecodeRequestBody {
			if err = request.Decode(w, r, &in); err != nil {
				encodeErr := response.JSON(w, 400, map[string] any {
					"code": 400,
					"messsage": err.Error(),
				})
	
				if encodeErr != nil {
					log.Printf("presentation.Post: %v", err)
				}
	
				return
			}

			defer r.Body.Close()
		}

		out := handle(r.Context(), in, commonInput)
		
		if opts.SetSessionCookie {
			setSessionCookie(w, out.CommonRes().SessionID, out.CommonRes().SessionMaxAge)
		}

		if err = response.JSON(w, out.CommonRes().StatusCode, out); err != nil {
			log.Printf("presentation.Post: %v", err)
			return
		}
	}
}

func Delete[outType b.CommonOutput](handle DefaultHandler[outType], opts Opts) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var commonInput = b.CommonInput{
			URLParam: opts.URLParams,
			QueryParam: opts.QueryParams,
		}

		var err error

		if opts.GetSessionCookie {
			commonInput.SessionID, err = getSessionCookie(r)
			if err != nil {
				encodeErr := response.JSON(w, 400, map[string]any{
					"code": 400,
					"messsage": err.Error(),
				})
	
				if encodeErr != nil {
					log.Printf("presentation.Get: %v", err)
				}
	
				return
			}
		}

		out := handle(r.Context(), commonInput)

		if opts.SetSessionCookie {
			setSessionCookie(w, out.CommonRes().SessionID, out.CommonRes().SessionMaxAge)
		}

		if err := response.JSON(w, out.CommonRes().StatusCode, &out); err != nil {
			log.Printf("presentation.Get: %v", err)
			return
		}
	}
}