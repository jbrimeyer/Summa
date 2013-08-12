package summa

import (
	"database/sql"
	"encoding/json"
	"errors"
	_ "go-sqlite3"
	"io/ioutil"
	"net/http"
)

const (
	API_STATUS_OK    = "ok"
	API_STATUS_ERROR = "error"
	API_STATUS_FATAL = "fatal"
	INTERNAL_ERROR   = "Internal application error"
)

type fatalError struct {
	s   string
	err error
}

type authError struct {
	err error
}

type apiError struct {
	err error
}

type permissionError struct {
	err error
}

type apiRequest struct {
	Username string      `json:"username"`
	Password string      `json:"password"`
	Token    string      `json:"token"`
	Data     interface{} `json:"data"`
}

type apiResponse struct {
	Status string                 `json:"status"`
	Error  string                 `json:"error,omitempty"`
	Token  string                 `json:"token,omitempty"`
	Data   map[string]interface{} `json:"data,omitempty"`
}

type apiHandlerFunc func(r *apiResponse, db *sql.DB) error

type apiHandler struct {
	handlerFunc   apiHandlerFunc
	isAuthHandler bool
}

var (
	errInvalidEndpoint    = errors.New("Invalid API endpoint")
	errInvalidMethod      = errors.New("Invalid HTTP request method")
	errMalformedJson      = errors.New("Malformed JSON request data")
	errInvalidCredentials = errors.New("Invalid authorization credentials")
	errInvalidSession     = errors.New("Invalid or expired session token")

	apiEndpoints = map[string]apiHandler{
		"/api/auth/signin": apiHandler{
			isAuthHandler: true,
		},
		"/api/auth/signout": apiHandler{
			handlerFunc: apiAuthSignout,
		},
		"/api/profile": apiHandler{
			handlerFunc: apiProfile,
		},
		"/api/profile/update": apiHandler{
			handlerFunc: apiProfileUpdate,
		},
		"/api/snippet": apiHandler{
			handlerFunc: apiSnippet,
		},
		"/api/snippet/create": apiHandler{
			handlerFunc: apiSnippetCreate,
		},
		"/api/snippet/update": apiHandler{
			handlerFunc: apiSnippetUpdate,
		},
		"/api/snippet/delete": apiHandler{
			handlerFunc: apiSnippetDelete,
		},
		"/api/search": apiHandler{
			handlerFunc: apiSearch,
		},
		"/api/unread": apiHandler{
			handlerFunc: apiUnread,
		},
	}
)

func (e *fatalError) Error() string {
	return INTERNAL_ERROR
}

func (e *authError) Error() string {
	return e.err.Error()
}

func (e *apiError) Error() string {
	return e.err.Error()
}

func (e *permissionError) Error() string {
	return e.err.Error()
}

func handleApiRequest(w http.ResponseWriter, req *http.Request) {
	header := w.Header()
	header.Set("Server", "Summa/1.0.0")
	header.Set("Content-Type", "application/json; charset=UTF-8")

	var resp apiResponse
	err := generateApiResponse(req, &resp)

	//
	// 200 OK
	//		Everthing is peachy, good request sent by the
	//		client and appropriate response returned
	// 400 Bad Request
	//		Client sent a bad request.  Possible malformed
	//		JSON or missing fields.
	// 401 Unauthorized
	//		User did not provide correct login credentials
	//		or their authorization token is invalid/expired
	// 403 Forbidden
	//		When the user tries to do something they are
	//		not allowed to do.  For example, edit a snippet
	//    that they do not own
	// 404 Not Found
	//		When the client tries to access an endpoint that
	//		does not exist
	// 500 Internal Server Error
	//		Some type of internal error in Summa :(
	//

	if err != nil {
		resp.Data = nil
		resp.Error = err.Error()

		switch err.(type) {
		case *fatalError:
			errLog.Printf("%s: %s", err.(*fatalError).s, err.(*fatalError).err)
			resp.Status = API_STATUS_FATAL

		case *apiError, *authError:
			resp.Status = API_STATUS_ERROR
			resp.Data = nil

		case error:
			errLog.Printf("%s", err)
			resp.Status = API_STATUS_FATAL
		}
	} else {
		resp.Status = API_STATUS_OK
	}

	b, err := json.Marshal(resp)
	if err != nil {
		errLog.Printf("json.Marshal() failed: %s", err)
		b, _ = json.Marshal(map[string]string{
			"status": API_STATUS_FATAL,
			"error":  INTERNAL_ERROR,
		})
	}

	w.Write(b)
}

func generateApiResponse(httpReq *http.Request, apiResp *apiResponse) error {
	if httpReq.Method != "POST" {
		return &apiError{errInvalidMethod}
	}

	handler, ok := apiEndpoints[httpReq.URL.Path]
	if !ok {
		return &apiError{errInvalidEndpoint}
	}

	body, err := ioutil.ReadAll(httpReq.Body)
	if err != nil {
		return &fatalError{"Could not read request body", err}
	}

	var apiReq apiRequest
	err = json.Unmarshal(body, &apiReq)
	if err != nil {
		return &apiError{errMalformedJson}
	}

	db, err := sql.Open("sqlite3", config.DBFile())
	if err != nil {
		return &fatalError{"Could not open database", err}
	}
	defer db.Close()

	authenticated := false
	if handler.isAuthHandler {
		// TODO: External authentication

		if !authenticated {
			return &authError{errInvalidCredentials}
		}

		token, err := createSession(db, apiReq.Username)
		if err != nil {
			return &fatalError{"Could not create session", err}
		}
		apiResp.Token = token

		return nil
	} else {
		authenticated, err = isValidSession(db, apiReq.Username, apiReq.Token)
		if err != nil {
			return &fatalError{"Could not check for valid session", err}
		}

		if !authenticated {
			return &authError{errInvalidSession}
		}

		err = handler.handlerFunc(apiResp, db)
		return err
	}
}

func apiAuthSignin(r *apiResponse, db *sql.DB) error {
	return nil
}

func apiAuthSignout(r *apiResponse, db *sql.DB) error {
	return nil
}

func apiProfile(r *apiResponse, db *sql.DB) error {
	return nil
}

func apiProfileUpdate(r *apiResponse, db *sql.DB) error {
	return nil
}

func apiSnippet(r *apiResponse, db *sql.DB) error {
	snippet, err := SnippetFetch(db, 1)
	if err != nil {
		return err
	}

	r.Data = make(map[string]interface{})
	r.Data["snippet"] = snippet

	return nil
}

func apiSnippetCreate(r *apiResponse, db *sql.DB) error {
	return nil
}

func apiSnippetUpdate(r *apiResponse, db *sql.DB) error {
	return nil
}

func apiSnippetDelete(r *apiResponse, db *sql.DB) error {
	return nil
}

func apiSearch(r *apiResponse, db *sql.DB) error {
	return nil
}

func apiUnread(r *apiResponse, db *sql.DB) error {
	return nil
}
