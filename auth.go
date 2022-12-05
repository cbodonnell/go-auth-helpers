package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Authenticate(w http.ResponseWriter, r *http.Request, endpoint string) (map[string]interface{}, error) {
	var auth map[string]interface{}
	client := &http.Client{}
	authReq, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return auth, err
	}
	for _, cookie := range r.Cookies() {
		authReq.AddCookie(cookie)
	}
	authReq.Header.Set("Accept", "application/json")
	authResp, err := client.Do(authReq)
	if err != nil {
		return auth, err
	}
	defer authResp.Body.Close()
	if authResp.StatusCode != http.StatusOK {
		return auth, fmt.Errorf("received status code %d from auth endpoint", authResp.StatusCode)
	}
	for _, cookie := range authResp.Cookies() {
		http.SetCookie(w, cookie)
		r.AddCookie(cookie)
	}
	err = json.NewDecoder(authResp.Body).Decode(&auth)
	if err != nil {
		return auth, err
	}
	return auth, nil
}
