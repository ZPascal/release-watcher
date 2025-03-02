package providers

import (
	"io/ioutil"
	"net/http"
	"testing"

	"gopkg.in/jarcoal/httpmock.v1"
)

func TestFetchJetBrainsReleases(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	testdata, err := ioutil.ReadFile("../testdata/jetbrains_releases.json")
	if err != nil {
		t.Errorf("Failed to load test data: %s", err)
	}

	httpmock.RegisterResponder(
		"GET", "https://data.services.jetbrains.com/products?code=GO",
		httpmock.NewStringResponder(200, string(testdata)),
	)

	provider := JetBrainsProvider{
		client: &http.Client{},
	}

	releases, err := provider.FetchReleases(&JetBrainsProject{Name: "go", Alias: "GoLand"})
	if err != nil {
		t.Error("Failed:", err)
	}

	if len(releases) != 22 {
		t.Error("Unexpected number of releases:", len(releases))
	}

	if releases[0].Name != "2018.1 (181.4203.567)" {
		t.Error("Unexpected release:", releases[0].Name)
	}
	if releases[1].Name != "2018.1 (181.4203.544 eap)" {
		t.Error("Unexpected release:", releases[1].Name)
	}

	if releases[0].Date.Year() != 2018 || releases[0].Date.Month() != 3 || releases[0].Date.Day() != 29 {
		t.Errorf("Unexpected date: %s", releases[0].Date.String())
	}

	if releases[0].URL != "https://download.jetbrains.com/go/goland-2018.1.tar.gz" {
		t.Error("Unexpected URL:", releases[0].URL)
	}
}
