package launch

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/antonzhukov/spacetrouble/internal/entity"
)

type spaceX struct {
	client *http.Client
	url    string
}

func NewSpaceX(client *http.Client, url string) *spaceX {
	return &spaceX{client: client, url: url}
}

func (s *spaceX) GetLaunches() (entity.Launches, error) {
	req, err := http.NewRequest("GET", s.url, nil)
	if err != nil {
		return nil, err
	}

	res, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected HTTP status: %d", res.StatusCode)
	}

	var launches entity.Launches
	err = json.NewDecoder(res.Body).Decode(&launches)
	if err != nil {
		return nil, err
	}

	return launches, nil
}
