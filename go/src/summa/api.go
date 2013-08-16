package summa

import (
	"database/sql"
	"encoding/json"
	_ "go-sqlite3"
	"io/ioutil"
	"net/http"
)

type apiRequestData map[string]interface{}

type apiRequest struct {
	Username string         `json:"username"`
	User     *User          `json:"-"`
	Password string         `json:"password"`
	Token    string         `json:"token"`
	Data     apiRequestData `json:"data"`
}

type apiResponseData map[string]interface{}

type apiResponse struct {
	Status int             `json:"-"`
	Error  string          `json:"error,omitempty"`
	Token  string          `json:"token,omitempty"`
	Data   apiResponseData `json:"data,omitempty"`
}

type apiHandlerFunc func(db *sql.DB, req apiRequest, resp apiResponseData) apiError

type apiHandler struct {
	handlerFunc   apiHandlerFunc
	isAuthHandler bool
}

var (
	apiAuthEndpoint = "/api/auth/signin"

	apiEndpoints = map[string]apiHandlerFunc{
		"/api/auth/signout":    apiAuthSignout,
		"/api/profile":         apiProfile,
		"/api/profile/update":  apiProfileUpdate,
		"/api/snippet":         apiSnippet,
		"/api/snippet/create":  apiSnippetCreate,
		"/api/snippet/update":  apiSnippetUpdate,
		"/api/snippet/delete":  apiSnippetDelete,
		"/api/comment/create":  apiCommentCreate,
		"/api/comment/update":  apiCommentUpdate,
		"/api/comment/delete":  apiCommentDelete,
		"/api/snippets/search": apiSnippetsSearch,
		"/api/snippets/unread": apiSnippetsUnread,
	}
)

func handleApiRequest(w http.ResponseWriter, req *http.Request) {
	header := w.Header()
	header.Set("Server", "Summa/1.0.0")
	header.Set("Content-Type", "application/json; charset=UTF-8")

	var resp apiResponse
	resp.Status = http.StatusOK

	apiErr := generateApiResponse(req, &resp)

	if apiErr != nil {
		resp.Status = apiErr.Code()
		resp.Error = apiErr.Error()
		resp.Data = apiErr.Data()

		switch apiErr.(type) {
		case *internalServerError:
			ise := apiErr.(*internalServerError)
			errLog.Printf("%s: %s", ise.s, ise.err)
			break
		}
	}

	b, err := json.Marshal(resp)
	if err != nil {
		errLog.Printf("json.Marshal() failed: %s", err)
		resp.Status = http.StatusInternalServerError
		b = []byte(INTERNAL_ERROR)
	}

	w.WriteHeader(resp.Status)
	w.Write(b)
}

func generateApiResponse(httpReq *http.Request, apiResp *apiResponse) apiError {
	if httpReq.Method != "POST" {
		return &methodNotAllowedError{}
	}

	isAuthHandler := httpReq.URL.Path == apiAuthEndpoint
	handler, ok := apiEndpoints[httpReq.URL.Path]
	if !ok && !isAuthHandler {
		return &badRequestError{"Unrecognized API endpoint"}
	}

	body, err := ioutil.ReadAll(httpReq.Body)
	if err != nil {
		return &internalServerError{"Could not read request body", err}
	}

	var apiReq apiRequest
	err = json.Unmarshal(body, &apiReq)
	if err != nil {
		return &badRequestError{"Malformed JSON"}
	}

	db, err := sql.Open("sqlite3", config.DBFile())
	if err != nil {
		return &internalServerError{"Could not open database", err}
	}
	defer db.Close()

	if isAuthHandler {
		apiReq.User, err = config.AuthProvider(apiReq.Username, apiReq.Password)
		if err != nil {
			return &internalServerError{"Could not authenticate user", err}
		}

		if apiReq.User == nil {
			return &unauthorizedError{"Invalid authentication credentials"}
		}

		exists, err := userExists(db, apiReq.User.Username)
		if err != nil {
			return &internalServerError{"Could not create session", err}
		}

		if !exists {
			err := userCreate(db, apiReq.User)
			if err != nil {
				return &internalServerError{"Could not create session", err}
			}
		}

		token, err := sessionCreate(db, apiReq.Username)
		if err != nil {
			return &internalServerError{"Could not create session", err}
		}

		apiResp.Token = token
		apiResp.Data = make(map[string]interface{})
		apiResp.Data["user"] = apiReq.User
		apiResp.Data["needEmail"] = apiReq.User.Email == ""

		return nil
	} else {
		apiReq.User, err = userFetch(db, apiReq.Username)
		if err != nil {
			return &internalServerError{"Could not fetch user", err}
		}

		authenticated, err := sessionIsValid(db, apiReq.Username, apiReq.Token)
		if err != nil {
			return &internalServerError{"Could not check for valid session", err}
		}

		if apiReq.User == nil || !authenticated {
			return &unauthorizedError{"Invalid or expired authentication session"}
		}

		if apiReq.Data == nil {
			apiReq.Data = make(map[string]interface{})
		}

		apiResp.Data = make(map[string]interface{})
		apiResp.Data["needEmail"] = apiReq.User.Email == ""

		apiErr := handler(db, apiReq, apiResp.Data)
		return apiErr
	}
}

func apiAuthSignout(db *sql.DB, req apiRequest, resp apiResponseData) apiError {
	err := sessionRemove(db, req.Username, req.Token)
	if err != nil {
		return &internalServerError{"Could not remove authentication session", err}
	}
	return nil
}

func apiProfile(db *sql.DB, req apiRequest, resp apiResponseData) apiError {
	// TODO: apiProfile()
	return nil
}

func apiProfileUpdate(db *sql.DB, req apiRequest, resp apiResponseData) apiError {
	// TODO: apiProfileUpdate()
	return nil
}
