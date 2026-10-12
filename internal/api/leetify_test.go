package api

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestLeetifyProfile_Success(t *testing.T) {
	t.Setenv("LEETIFY_AUTH_TOKEN", "mocked-token")
	useLeetifyServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/profile/id/testid" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		var profile Profile
		profile.Meta.Name = "TestUser"
		json.NewEncoder(w).Encode(profile)
	})

	profile, err := LeetifyProfile("testid")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if profile.Meta.Name != "TestUser" {
		t.Errorf("expected profile.Meta.Name = TestUser, got %s", profile.Meta.Name)
	}
}

func TestLeetifyProfile_SendsAuthHeaderOnce(t *testing.T) {
	t.Setenv("LEETIFY_AUTH_TOKEN", "mocked-token")
	var got []string
	useLeetifyServer(t, func(w http.ResponseWriter, r *http.Request) {
		got = r.Header.Values("Authorization")
		json.NewEncoder(w).Encode(Profile{})
	})

	if _, err := LeetifyProfile("testid"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(got) != 1 || got[0] != "Bearer mocked-token" {
		t.Errorf("Authorization headers = %q, want one bearer token", got)
	}
}

func TestLeetifyProfile_MalformedResponse(t *testing.T) {
	useLeetifyServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<html>maintenance</html>"))
	})

	if _, err := LeetifyProfile("testid"); err == nil {
		t.Fatal("expected an error, got none")
	}
}

func TestLeetifyProfile_ErrorStatuses(t *testing.T) {
	for _, status := range []int{http.StatusUnauthorized, http.StatusNotFound, http.StatusInternalServerError} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			t.Setenv("LEETIFY_AUTH_TOKEN", "mocked-token")
			useLeetifyServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(status)
			})

			if _, err := LeetifyProfile("testid"); err == nil {
				t.Fatal("expected an error, got none")
			}
		})
	}
}
