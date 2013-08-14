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
	Status      int             `json:"-"`
	Username    string          `json:"username,omitempty"`
	DisplayName string          `json:"displayName,omitempty"`
	HasEmail    bool            `json:"hasEmail,omitempty"`
	Error       string          `json:"error,omitempty"`
	Token       string          `json:"token,omitempty"`
	Data        apiResponseData `json:"data,omitempty"`
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

	if handler.isAuthHandler {
		u, err := config.AuthProvider(apiReq.Username, apiReq.Password)
		if err != nil {
			return &internalServerError{"Could not authenticate user", err}
		}

		if u == nil {
			return &unauthorizedError{"Invalid authentication credentials"}
		}

		exists, err := userExists(db, u.Username)
		if err != nil {
			return &internalServerError{"Could not create session", err}
		}

		if !exists {
			err := userCreate(db, u)
			if err != nil {
				return &internalServerError{"Could not create session", err}
			}
		}

		token, err := sessionCreate(db, apiReq.Username)
		if err != nil {
			return &internalServerError{"Could not create session", err}
		}

		apiResp.DisplayName = u.DisplayName
		apiResp.Username = u.Username
		apiResp.HasEmail = u.Email != ""
		apiResp.Token = token

		return nil
	} else {
		u, err := userFetch(db, apiReq.Username)
		if err != nil {
			return &internalServerError{"Could not fetch user", err}
		}

		authenticated, err := sessionIsValid(db, apiReq.Username, apiReq.Token)
		if err != nil {
			return &internalServerError{"Could not check for valid session", err}
		}

		if u == nil || !authenticated {
			return &unauthorizedError{"Invalid or expired authentication session"}
		}

		if apiReq.Data == nil {
			apiReq.Data = make(map[string]interface{})
		}

		apiResp.DisplayName = u.DisplayName
		apiResp.Username = u.Username
		apiResp.HasEmail = u.Email != ""

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
	id, ok := req["id"].(string)

	if !ok {
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
		filenames := make(map[string]bool)
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

				lcFilename := strings.ToLower(fields["filename"])
				_, ok := filenames[lcFilename]
				if !fileRegex.MatchString(fields["filename"]) || ok {
					return &conflictError{apiResponseData{"field": fmt.Sprintf("file[%d].filename", i)}}
				}

				filenames[lcFilename] = true

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

	snip.Files = files
	snip.Username = req["_username"].(string)

	id, err := snippetCreate(db, &snip)
	if err != nil {
		return &internalServerError{"Could not create snippet", err}
	}

	resp["id"] = id

	return nil
}

func apiSnippetUpdate(db *sql.DB, req apiRequestData, resp apiResponseData) apiError {
	// TODO: apiSnippetUpdate
	return nil
}

func apiSnippetDelete(db *sql.DB, req apiRequestData, resp apiResponseData) apiError {
	id, ok := req["id"].(string)

	if !ok {
		return &badRequestError{"The 'id' field must be a string"}
	}

	owned, err := snippetIsOwnedBy(db, id, req["_username"].(string))
	if err != nil {
		return &internalServerError{"Could not check snippet ownership", err}
	}

	if !owned {
		return &forbiddenError{"You do not have permission to delete this snippet"}
	}

	err = snippetDelete(db, id)
	if err != nil {
		return &internalServerError{"Could not delete snippet", err}
	}

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
