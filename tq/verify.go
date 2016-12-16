package tq

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/git-lfs/git-lfs/lfsapi"
)

func verifyUpload(c *lfsapi.Client, t *Transfer) error {
	action, err := t.Actions.Get("verify")
	if err != nil {
		if IsActionMissingError(err) {
			return nil
		}
		return err
	}

	by, err := json.Marshal(struct {
		Oid  string `json:"oid"`
		Size int64  `json:"size"`
	}{Oid: t.Oid, Size: t.Size})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", action.Href, bytes.NewReader(by))
	if err != nil {
		return err
	}

	for key, value := range action.Header {
		req.Header.Set(key, value)
	}
	req.Header.Set("Content-Type", "application/vnd.git-lfs+json")

	res, err := c.Do(req)
	if err != nil {
		return err
	}

	return res.Body.Close()
}
