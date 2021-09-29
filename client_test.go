package wemessage

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var MockClient = NewClient("mockCorpID", "mockCorpSecret")

func TestClient_RenewAccessToken(t *testing.T) {

	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			sampleResponse := `{"errcode": 0,"errmsg": "ok","access_token": "mockAccessToken","expires_in": 7200}`
			_, _ = w.Write([]byte(sampleResponse))
		}),
	)

	MockClient.BaseURL = ts.URL

	if err := MockClient.RenewAccessToken(); err != nil {
		t.Fatal(err)
	}
	t.Log(MockClient)
}
