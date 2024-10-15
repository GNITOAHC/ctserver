package router

import (
	"bytes"
	// "ctserver/internal/helper"
	"encoding/json"
	"io"
	"net/http"
)

type UrlReq struct {
	Username string `json:"username"` // User's mail
	Url      string `json:"url"`      // Url to shorten
	Name     string `json:"name"`     // Shortened name of the url, empty for random
	Desc     string `json:"desc"`     // Short description of the url
	Ancestor string `json:"ancestor"` // Ancestor id, "/" for root
}

func (rr *Router) NewUrl(w http.ResponseWriter, r *http.Request) {
	var u UrlReq
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// shortened, err := rr.helper.ShortenUrl(helper.Url{
	// 	Username: u.Username,
	// 	Url:      u.Url,
	// 	Desc:     u.Desc,
	// 	Ancestor: u.Ancestor,
	// })
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// json.NewEncoder(w).Encode(map[string]string{
	// 	"shortened": shortened,
	// })
	return
}

func (rr *Router) NewFile(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	r.ParseMultipartForm(10 << 20)
	var buf bytes.Buffer
	file, _, err := r.FormFile("file")
	defer file.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	io.Copy(&buf, file)
	// rr.helper.UploadFile(helper.File{
	// 	Mail:     r.FormValue("mail"),
	// 	Name:     r.FormValue("filename"),
	// 	Ancestor: r.FormValue("ancestor"),
	// 	File:     buf.Bytes(),
	// })
}
