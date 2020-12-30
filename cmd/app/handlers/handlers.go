package handlers

import (

	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Xusrav/GoAuth2.0/pkg/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)
var (
	googleOauthConfig *oauth2.Config
	outhStateString = "pseudo-random"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func init() {
	redirectUrl := "http://"+config.Host+":"+config.Port+"/redirect"
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  redirectUrl,
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}

func (h *Handler)HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL(outhStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handler)HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	content, err := h.getUserInfo(r.FormValue("state"), r.FormValue("code"))
	if err != nil {
		fmt.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	fmt.Fprintf(w, "Content: %s\n", content)
}
func (h *Handler)getUserInfo(state string, code string) ([]byte, error) {
	if state != outhStateString {
		return nil, fmt.Errorf("invalid oauth state")
	}
	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}
	return contents, nil
}

func (h *Handler)HandleMain(w http.ResponseWriter, r *http.Request) {
	var htmlIndex = `<html>
<body>
  <a href="/login">Google Log In</a>
</body>
</html>`
	fmt.Fprintf(w, htmlIndex)
}