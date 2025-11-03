package swapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Name      string   `json:"name"`
	Gravity   string   `json:"gravity"`
	Residents []string `json:"residents"`
}

type client struct {
	httpClient *http.Client
}

func NewClient(httpClient *http.Client) *client {
	return &client{
		httpClient: httpClient,
	}
}

func (c *client) GetPlanet(planets int64) (Response, error) {
	res, err := c.httpClient.Get(fmt.Sprintf("https://swapi.dev/api/planets/%d", planets))
	if err != nil {
		return Response{}, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return Response{}, fmt.Errorf("status code %d", res.StatusCode)
	}

	var r Response
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return Response{}, err
	}

	return r, nil
}
