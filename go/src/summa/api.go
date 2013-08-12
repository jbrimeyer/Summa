package summa

import (
	"database/sql"
	"encoding/json"
	_ "go-sqlite3"
	"io/ioutil"
	"net/http"
)

type apiRequest struct {
	Username string      `json:"username"`
	Password string      `json:"password"`
	Token    string      `json:"token"`
	Data     interface{} `json:"data"`
}

type apiResponse struct {
	Status int                    `json:"-"`
	Error  string                 `json:"error,omitempty"`
	Token  string                 `json:"token,omitempty"`
	Data   map[string]interface{} `json:"data,omitempty"`
}

type apiHandlerFunc func(db *sql.DB, req interface{}, resp *apiResponse) apiError

type apiHandler struct {
	handlerFunc   apiHandlerFunc
	isAuthHandler bool
}

var (
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
		resp.Data = nil

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
		// TODO: Generate internalServerError
	}

	w.WriteHeader(resp.Status)
	w.Write(b)
}

func generateApiResponse(httpReq *http.Request, apiResp *apiResponse) apiError {
	if httpReq.Method != "POST" {
		return &badRequestError{"Invalid request method"}
	}

	handler, ok := apiEndpoints[httpReq.URL.Path]
	if !ok {
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

	authenticated := false
	if handler.isAuthHandler {
		// TODO: External authentication

		if !authenticated {
			return &unauthorizedError{"Invalid authentication credentials"}
		}

		token, err := createSession(db, apiReq.Username)
		if err != nil {
			return &internalServerError{"Could not create session", err}
		}
		apiResp.Token = token

		return nil
	} else {
		authenticated, err = isValidSession(db, apiReq.Username, apiReq.Token)
		if err != nil {
			return &internalServerError{"Could not check for valid session", err}
		}

		if !authenticated {
			return &unauthorizedError{"Invalid or expired authentication session"}
		}

		apiErr := handler.handlerFunc(db, apiReq.Data, apiResp)
		return apiErr
	}
}

func apiAuthSignin(db *sql.DB, req interface{}, resp *apiResponse) apiError {
	return nil
}

func apiAuthSignout(db *sql.DB, req interface{}, resp *apiResponse) apiError {
	return nil
}

func apiProfile(db *sql.DB, req interface{}, resp *apiResponse) apiError {
	return nil
}

func apiProfileUpdate(db *sql.DB, req interface{}, resp *apiResponse) apiError {
	return nil
}

func apiSnippet(db *sql.DB, req interface{}, resp *apiResponse) apiError {
	snippet, err := SnippetFetch(db, 1)
	if err != nil {
		return &internalServerError{"Could not fetch snippet", err}
	}

	if snippet == nil {
		return &notFoundError{"No such snippet"}
	}

	resp.Data = make(map[string]interface{})
	resp.Data["snippet"] = snippet

	return nil
}

func apiSnippetCreate(db *sql.DB, req interface{}, resp *apiResponse) apiError {
	return nil
}

func apiSnippetUpdate(db *sql.DB, req interface{}, resp *apiResponse) apiError {
	return nil
}

func apiSnippetDelete(db *sql.DB, req interface{}, resp *apiResponse) apiError {
	return nil
}

func apiSearch(db *sql.DB, req interface{}, resp *apiResponse) apiError {
	return nil
}

func apiUnread(db *sql.DB, req interface{}, resp *apiResponse) apiError {
	return nil
}
