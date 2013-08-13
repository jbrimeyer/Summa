package summa

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "go-sqlite3"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

type apiRequestData map[string]interface{}

type apiRequest struct {
	Username string         `json:"username"`
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

type apiHandlerFunc func(db *sql.DB, req apiRequestData, resp apiResponseData) apiError

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
		// TODO: Generate internalServerError
	}

	w.WriteHeader(resp.Status)
	w.Write(b)
}

func generateApiResponse(httpReq *http.Request, apiResp *apiResponse) apiError {
	if httpReq.Method != "POST" {
		return &methodNotAllowedError{}
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
		// TODO: If user is new, add them to the user database table
		// TODO: If user is new, indicate that in the JSON response so the client can prompt for e-mail address

		if !authenticated {
			return &unauthorizedError{"Invalid authentication credentials"}
		}

		token, err := sessionCreate(db, apiReq.Username)
		if err != nil {
			return &internalServerError{"Could not create session", err}
		}
		apiResp.Token = token

		return nil
	} else {
		authenticated, err = sessionIsValid(db, apiReq.Username, apiReq.Token)
		if err != nil {
			return &internalServerError{"Could not check for valid session", err}
		}

		if !authenticated {
			return &unauthorizedError{"Invalid or expired authentication session"}
		}

		if apiReq.Data == nil {
			apiReq.Data = make(map[string]interface{})
		}

		apiReq.Data["_username"] = apiReq.Username
		apiReq.Data["_token"] = apiReq.Token
		apiResp.Data = make(map[string]interface{})
		apiErr := handler.handlerFunc(db, apiReq.Data, apiResp.Data)
		return apiErr
	}
}

func apiAuthSignout(db *sql.DB, req apiRequestData, resp apiResponseData) apiError {
	err := sessionRemove(db, req["_username"].(string), req["_token"].(string))
	if err != nil {
		return &internalServerError{"Could not remove authentication session", err}
	}
	return nil
}

func apiProfile(db *sql.DB, req apiRequestData, resp apiResponseData) apiError {
	// TODO: apiProfile()
	return nil
}

func apiProfileUpdate(db *sql.DB, req apiRequestData, resp apiResponseData) apiError {
	// TODO: apiProfileUpdate()
	return nil
}

func apiSnippet(db *sql.DB, req apiRequestData, resp apiResponseData) apiError {
	var id int64
	var err error
	switch req["id"].(type) {
	case string:
		id, err = FromBase36(req["id"].(string))
		if err != nil {
			return &badRequestError{"Invalid snippet id"}
		}
	default:
		return &badRequestError{"The 'id' field must be a string"}
	}

	snippet, err := snippetFetch(db, id)
	if err != nil {
		return &internalServerError{"Could not fetch snippet", err}
	}

	if snippet == nil {
		return &notFoundError{"No such snippet"}
	}

	resp["snippet"] = snippet

	return nil
}

func apiSnippetCreate(db *sql.DB, req apiRequestData, resp apiResponseData) apiError {
	reqForFiles := []string{"filename", "language", "contents"}
	var snip snippet

	snip.Description, _ = req["description"].(string)
	if snip.Description == "" {
		return &conflictError{apiResponseData{"field": "description"}}
	}

	fileRegex := regexp.MustCompile("(?i)^[a-z0-9_.-]+$")
	var files snippetFiles

	switch req["files"].(type) {
	case []interface{}:
		for i, v := range req["files"].([]interface{}) {
			switch v.(type) {
			case map[string]interface{}:
				vmap := v.(map[string]interface{})
				fields := make(map[string]string)
				var file snippetFile

				for _, required := range reqForFiles {
					strVal, ok := vmap[required].(string)
					if !ok || strings.TrimSpace(strVal) == "" {
						return &conflictError{apiResponseData{"field": fmt.Sprintf("file[%d].%s", i, required)}}
					}
					fields[required] = strings.TrimSpace(strVal)
				}

				if !fileRegex.MatchString(fields["filename"]) {
					return &conflictError{apiResponseData{"field": fmt.Sprintf("file[%d].filename", i)}}
				}

				file.Filename = fields["filename"]
				file.Language = fields["language"]
				file.Contents = fields["contents"]

				files = append(files, file)
			default:
				return &badRequestError{"'files' field is malformed"}
			}
		}
	default:
		return &conflictError{apiResponseData{"field": "files"}}
	}

	snip.Files = &files

	// TODO: Start database transaction
	// TODO: Insert snippet into database
	// TODO: Create git repository
	// TODO: Any problems? roll back the transaction and cleanup filesystem
	// TODO: All ok? return the snippet ID in resp

	return nil
}

func apiSnippetUpdate(db *sql.DB, req apiRequestData, resp apiResponseData) apiError {
	// TODO: apiSnippetUpdate
	return nil
}

func apiSnippetDelete(db *sql.DB, req apiRequestData, resp apiResponseData) apiError {
	// TODO: apiSnippetDelete
	return nil
}

func apiSearch(db *sql.DB, req apiRequestData, resp apiResponseData) apiError {
	// TODO: apiSearch()
	return nil
}

func apiUnread(db *sql.DB, req apiRequestData, resp apiResponseData) apiError {
	snippets, err := snippetsUnread(db, req["_username"].(string))
	if err != nil {
		return &internalServerError{"Could not fetch snippets", err}
	}

	resp["snippets"] = snippets

	return nil
}
