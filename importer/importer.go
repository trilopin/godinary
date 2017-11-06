package importer

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/trilopin/godinary/storage"
)

type Importer interface {
	Import(sd storage.Driver)
}

type CloudinaryImporter struct {
	UserSpace string
	APIKey    string
	APISecret string
}

type CloudinaryResult struct {
	PublicID string `json:"public_id"`
	Format   string `json:"format"`
	URL      string `json:"url"`
}

type CloudinaryResponse struct {
	Resources  []CloudinaryResult `json:"resources"`
	NextCursor string             `json:"next_cursor"`
}

func (ci *CloudinaryImporter) NewRequest(cursor string) (*http.Request, error) {
	cURL, err := url.Parse(fmt.Sprintf("https://api.cloudinary.com/v1_1/%s/resources/search", ci.UserSpace))
	if err != nil {
		return nil, err
	}

	q := cURL.Query()
	q.Set("expression", "type:upload")
	q.Set("max_results", "100")
	if cursor != "" {
		q.Set("next_cursor", cursor)
	}
	cURL.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", cURL.String(), nil)
	req.SetBasicAuth(ci.APIKey, ci.APISecret)
	return req, err
}

func (cr *CloudinaryResponse) Upload(sd storage.Driver) {
	var wg sync.WaitGroup
	for _, r := range cr.Resources {
		wg.Add(1)
		go func(r CloudinaryResult, sd storage.Driver) {
			defer wg.Done()
			t1 := time.Now()
			name := fmt.Sprintf("%s.%s", r.PublicID, r.Format)

			ht := sha256.New()
			ht.Write([]byte(name))
			hash := hex.EncodeToString(ht.Sum(nil))

			resp, err := http.Get(r.URL)
			t2 := time.Now()
			if err != nil {
				log.Println("Failed to download", r.URL)
				return
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println("Failed to parse body", r.URL)
				return
			}

			reader, err := sd.NewReader(hash, "upload/")
			if err == nil {
				defer reader.Close()
			} else {
				sd.Write(body, hash, "upload/")
				fmt.Printf("\nget %s, write %s - %s ", t2.Sub(t1), time.Since(t2), r.URL)
			}

		}(r, sd)
	}
	t3 := time.Now()
	wg.Wait()
	fmt.Printf("\nwait %s\n", time.Since(t3))

}

func (ci *CloudinaryImporter) Import(sd storage.Driver) error {
	var cursor string
	var err error
	cursor = ""
	client := &http.Client{}
	for {
		req, err := ci.NewRequest(cursor)
		if err != nil {
			return fmt.Errorf("can not create request: %v", err)
		}
		resp, err := client.Do(req)

		if err != nil {
			return fmt.Errorf("can not download: %v", err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		cd := &CloudinaryResponse{}
		json.Unmarshal(body, cd)

		cd.Upload(sd)
		if cd.NextCursor != "" {
			fmt.Println(cd.NextCursor)
			cursor = cd.NextCursor
		} else {
			break
		}
	}
	return err
}
