package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Xusrav/GoAuth2.0/pkg/config"
	"github.com/imroc/req"
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


func (h *Handler)HandleSearch(w http.ResponseWriter, r *http.Request) {
	getParamFromRequestFormData(r)
	param := ""
	if by == "id" {
		if id == "" {
			w.WriteHeader(404)
			w.Write([]byte("Неправильный запорс"))
			return
		}
		param = "&i=" + id
	} else if by == "title" {
		if title == "" {
			w.WriteHeader(404)
			w.Write([]byte("Неправильный запорс"))
			return
		}
		param= "&t=" + title
	} else if by == "search" {
		if s == "" {
			w.WriteHeader(404)
			w.Write([]byte("Неправильный запорс"))
			return
		}
		param= "&s=" + s
	}

	if year != "" {
		param += "&y=" + year
	}
	if plot != "" {
		param += "&plot=" + plot
	}
	if typeData != "" {
		param += "&r=" + typeData
	}

	if typeMovie != "" {
		param += "&type=" + typeMovie
	}

	if page != "" {
		param += "&page=" + page
	}

	url := config.URLomdbApi+"?apikey=" + config.ApiKey

	log.Println(url+param)
	post, err := req.Get(url+param)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Ошибка сервера"))
		return
	}

	if typeData !="xml"{
		w.Header().Set("Content-type", "application/json")
	}else {
		w.Header().Set("Content-type", "application/xml")
	}
	w.Write(post.Bytes())
	return
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
	w.Header().Set("Content-type", "text/html")

	tpl, err := template.ParseFiles(
		filepath.Join("../pkg/templates/html.gohtml"),
		filepath.Join("../pkg/templates/base.gohtml"),
	)

	err = tpl.Execute(w, []byte("privet"))
	if err != nil {
		log.Println(err)
	}
}


var by, id, year, plot, typeData, s, page, title, typeMovie string

func getParamFromRequestFormData(r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	a := string(bytes)
	b := strings.ReplaceAll(a, "Content-Disposition: form-data; ", "")
	b = strings.ReplaceAll(b, "\r\n\r\n", "")
	c := strings.Split(b, "\r\n")

	for i := 1; i < len(c); i += 2 {
		text := strings.Split(c[i], `"`)
		if len(text) < 3 {
			continue
		}
		if text[1] == "by" {
			by = text[2]
		}
		if text[1] == "i" {
			id = text[2]
		}
		if text[1] == "t" {
			title = text[2]
		}
		if text[1] == "type" {
			typeMovie = text[2]
		}
		if text[1] == "y" {
			year = text[2]
		}
		if text[1] == "plot" {
			plot = text[2]
		}
		if text[1] == "r" {
			typeData = text[2]
		}
		if text[1] == "s" {
			s = text[2]
		}
		if text[1] == "page" {
			page = text[2]
		}
	}
}