package players

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joho/godotenv"
)

func TestWinRecordingAndRetrieving(t *testing.T) {
	envErr := godotenv.Load("../.env")
	if envErr != nil {
		t.Fatalf("error while loading env variables: %s", envErr)
	}

	conn, err := ConnectToDB()

	if err != nil {
		t.Fatalf("error while connecting to db: %s", err)
	}

	store := PostgresPlayerStore{Conn: conn}
	svr := PlayerServer{&store}
	prep := DBPrep{conn, t}

	player := "Pepper"
	prep.deletePlayer(player)
	for range 3 {
		svr.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	}

	res := httptest.NewRecorder()
	req := newGetScoreRequest(player)
	svr.ServeHTTP(res, req)

	assertStatus(t, res, http.StatusOK)
	assertResponseBody(t, res, "3")
}
