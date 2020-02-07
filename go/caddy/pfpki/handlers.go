package pfpki

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/inverse-inc/packetfence/go/log"
)

type (
	// Info struct
	Info struct {
		Status      int    `json:"status"`
		Password    string `json:"password"`
		Error       string `json:"error"`
		ContentType string
		Raw         []byte
		Entries     interface{} `json:"items"`
		NextCursor  int         `json:"nextCursor"`
		PrevCursor  int         `json:"prevCursor"`
		TotalCount  int         `json:"total_count"`
	}
)

var decoder = schema.NewDecoder()

// Create interface
type Create interface {
	new() error
	get() error
	revoke() error
}

func manageCA(pfpki *Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		var o CA

		manage(o, pfpki, res, req)
	})
}

func manageCert(pfpki *Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		var o Cert

		manage(o, pfpki, res, req)
	})
}

func manageProfile(pfpki *Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		var o Profile

		manage(o, pfpki, res, req)
	})
}

func manage(object interface{}, pfpki *Handler, res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	body, err := ioutil.ReadAll(req.Body)
	var Information Info
	// Set the default Content-Type
	Information.ContentType = "application/json; charset=UTF-8"

	if err != nil {
		panic(err)
	}

	switch v := object.(type) {
	case CA:

		switch {

		case len(regexp.MustCompile(`/pki/cas/search$`).FindStringIndex(req.URL.Path)) > 0:
			switch req.Method {
			case "POST":
				var params PostVars
				err := json.Unmarshal(body, &params)
				if err != nil {
					log.LoggerWContext(pfpki.Ctx).Info(err.Error())
					panic(err)
				}
				pagination := params.Sanitize(object)
				Information, err = v.search(pfpki, pagination)
			default:
				err = errors.New("Method not supported")
				log.LoggerWContext(pfpki.Ctx).Info("Method not supported")
			}

		case len(regexp.MustCompile(`/pki/cas`).FindStringIndex(req.URL.Path)) > 0:
			switch req.Method {
			case "GET":
				var params GetVars
				err := decoder.Decode(&params, req.URL.Query())
				if err != nil {
					log.LoggerWContext(pfpki.Ctx).Info(err.Error())
					panic(err)
				}
				pagination := params.Sanitize(object)
				Information, err = v.paginated(pfpki, pagination)
			case "POST":
				err = json.Unmarshal(body, &v)
				if err != nil {
					log.LoggerWContext(pfpki.Ctx).Info(err.Error())
					panic(err)
				}
				Information, err = v.new(pfpki)
			default:
				err = errors.New("Method not supported")
				log.LoggerWContext(pfpki.Ctx).Info("Method not supported")
			}

		case len(regexp.MustCompile(`/pki/ca/[0-9]+$`).FindStringIndex(req.URL.Path)) > 0:
			switch req.Method {
			case "GET":
				Information, err = v.getById(pfpki, vars)
			default:
				err = errors.New("Method not supported")
				log.LoggerWContext(pfpki.Ctx).Info("Method not supported")
			}

		default:
			err = errors.New("Path not supported")
			log.LoggerWContext(pfpki.Ctx).Info("Path not supported" + req.URL.Path)
		}

	case Profile:

		switch {

		case len(regexp.MustCompile(`/pki/profiles/search$`).FindStringIndex(req.URL.Path)) > 0:
			switch req.Method {
			case "POST":
				var params PostVars
				err := json.Unmarshal(body, &params)
				if err != nil {
					log.LoggerWContext(pfpki.Ctx).Info(err.Error())
					panic(err)
				}
				pagination := params.Sanitize(object)
				Information, err = v.search(pfpki, pagination)
			default:
				err = errors.New("Method not supported")
				log.LoggerWContext(pfpki.Ctx).Info("Method not supported")
			}

		case len(regexp.MustCompile(`/pki/profiles`).FindStringIndex(req.URL.Path)) > 0:
			switch req.Method {
			case "GET":
				var params GetVars
				err := decoder.Decode(&params, req.URL.Query())
				if err != nil {
					log.LoggerWContext(pfpki.Ctx).Info(err.Error())
					panic(err)
				}
				pagination := params.Sanitize(object)
				Information, err = v.paginated(pfpki, pagination)
			case "POST":
				err = json.Unmarshal(body, &v)
				if err != nil {
					log.LoggerWContext(pfpki.Ctx).Info(err.Error())
					panic(err)
				}
				Information, err = v.new(pfpki)
			default:
				err = errors.New("Method not supported")
				log.LoggerWContext(pfpki.Ctx).Info("Method not supported")
			}

		case len(regexp.MustCompile(`/pki/profile/[0-9]+$`).FindStringIndex(req.URL.Path)) > 0:
			switch req.Method {
			case "GET":
				Information, err = v.getById(pfpki, vars)
			default:
				err = errors.New("Method not supported")
				log.LoggerWContext(pfpki.Ctx).Info("Method not supported")
			}

		default:
			err = errors.New("Path not supported")
			log.LoggerWContext(pfpki.Ctx).Info("Path not supported")
		}

	case Cert:

		switch {

		case len(regexp.MustCompile(`/pki/certs/search$`).FindStringIndex(req.URL.Path)) > 0:
			switch req.Method {
			case "POST":
				var params PostVars
				err := json.Unmarshal(body, &params)
				if err != nil {
					log.LoggerWContext(pfpki.Ctx).Info(err.Error())
					panic(err)
				}
				pagination := params.Sanitize(object)
				Information, err = v.search(pfpki, pagination)
			default:
				err = errors.New("Method not supported")
				log.LoggerWContext(pfpki.Ctx).Info("Method not supported")
			}

		case len(regexp.MustCompile(`/pki/certs`).FindStringIndex(req.URL.Path)) > 0:
			switch req.Method {
			case "GET":
				var params GetVars
				err := decoder.Decode(&params, req.URL.Query())
				if err != nil {
					log.LoggerWContext(pfpki.Ctx).Info(err.Error())
					panic(err)
				}
				pagination := params.Sanitize(object)
				Information, err = v.paginated(pfpki, pagination)
			case "POST":
				err = json.Unmarshal(body, &v)
				if err != nil {
					log.LoggerWContext(pfpki.Ctx).Info(err.Error())
					panic(err)
				}
				Information, err = v.new(pfpki)
			default:
				err = errors.New("Method not supported")
				log.LoggerWContext(pfpki.Ctx).Info("Method not supported")
			}

		case len(regexp.MustCompile(`/pki/cert/[0-9]+$`).FindStringIndex(req.URL.Path)) > 0:
			switch req.Method {
			case "GET":
				Information, err = v.getById(pfpki, vars)
			default:
				err = errors.New("Method not supported")
				log.LoggerWContext(pfpki.Ctx).Info("Method not supported")
			}

		case len(regexp.MustCompile(`/pki/cert/[0-9]+/download/.*$`).FindStringIndex(req.URL.Path)) > 0:
			switch req.Method {
			case "GET":
				Information, err = v.download(pfpki, vars)
			default:
				err = errors.New("Method not supported")
				log.LoggerWContext(pfpki.Ctx).Info("Method not supported")
			}

		case len(regexp.MustCompile(`/pki/cert/[0-9]+/email$`).FindStringIndex(req.URL.Path)) > 0:
			switch req.Method {
			case "GET":
				Information, err = v.download(pfpki, vars)
			default:
				err = errors.New("Method not supported")
				log.LoggerWContext(pfpki.Ctx).Info("Method not supported")
			}

		case len(regexp.MustCompile(`/pki/cert/[0-9]+/[0-9]+$`).FindStringIndex(req.URL.Path)) > 0:
			switch req.Method {
			case "DELETE":
				Information, err = v.revoke(pfpki, vars)
			default:
				err = errors.New("Method not supported")
				log.LoggerWContext(pfpki.Ctx).Info("Method not supported")
			}

		default:
			err = errors.New("Path not supported")
			log.LoggerWContext(pfpki.Ctx).Info("Path not supported")
		}

	default:
		err = errors.New("Type not supported")
	}

	if err != nil {
		Information.Error = err.Error()
	}

	var result = &Information

	switch ContentType := Information.ContentType; ContentType {

	case "application/x-pkcs12":
		res.Header().Set("Content-Type", "application/x-pkcs12")
		res.Header().Set("Content-Disposition", "attachment; filename=certificate.p12")
		res.WriteHeader(http.StatusOK)
		io.Copy(res, bytes.NewReader(Information.Raw))

	default:
		res.Header().Set("Content-Type", "application/json; charset=UTF-8")
		res.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(res).Encode(result); err != nil {
			panic(err)
		}
	}
}

func manageOcsp(pfpki *Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		log.LoggerWContext(pfpki.Ctx).Info(fmt.Sprintf("Got %s request from %s", req.Method, req.RemoteAddr))
		if req.Header.Get("Content-Type") != "application/ocsp-request" {
			log.LoggerWContext(pfpki.Ctx).Info("Strict mode requires correct Content-Type header")
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		b := new(bytes.Buffer)
		switch req.Method {
		case "POST":
			b.ReadFrom(req.Body)
		case "GET":
			log.LoggerWContext(pfpki.Ctx).Info(req.URL.Path)
			gd, err := base64.StdEncoding.DecodeString(req.URL.Path[1:])
			if err != nil {
				log.LoggerWContext(pfpki.Ctx).Info(err.Error())
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			r := bytes.NewReader(gd)
			b.ReadFrom(r)
		default:
			log.LoggerWContext(pfpki.Ctx).Info("Unsupported request method")
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		oscp := Responder(pfpki)
		// parse request, verify, create response
		res.Header().Set("Content-Type", "application/ocsp-response")
		resp, err := oscp.verify(b.Bytes())
		if err != nil {
			log.LoggerWContext(pfpki.Ctx).Info(err.Error())
			// technically we should return an ocsp error response. but this is probably fine
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		log.LoggerWContext(pfpki.Ctx).Info("Writing response")
		res.Write(resp)
	})
}