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

// Importer is the interface to implement to imretrieve and save images
type Importer interface {
	Import(sd storage.Driver)
}

// CloudinaryImporter is the element to import uploaded images from claoudinary
type CloudinaryImporter struct {
	UserSpace string
	APIKey    string
	APISecret string
}

// CloudinaryResult is the struct for json item unmarshalling
type CloudinaryResult struct {
	PublicID string `json:"public_id"`
	Format   string `json:"format"`
	URL      string `json:"url"`
}

// CloudinaryResponse is the struct for json response unmarshalling
type CloudinaryResponse struct {
	Resources  []CloudinaryResult `json:"resources"`
	NextCursor string             `json:"next_cursor"`
}

// NewRequest creates a new HTTP Request preconfigured to retrieve images
// from cloudinary. Take care about, pagination, userspace and auth
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

// Upload save the images from the collection URLs into storage
// The hash is build based on filename
func (cr *CloudinaryResponse) Upload(sd storage.Driver) {
	var wg sync.WaitGroup
	for _, r := range cr.Resources {
		wg.Add(1)
		err := sd.Init() // force connection
		if err != nil {
			fmt.Printf("can't initalise storage: %v", err)
			return
		}

		go func(r CloudinaryResult, sd storage.Driver) {
			defer wg.Done()
			t1 := time.Now()
			name := fmt.Sprintf("%s.%s", r.PublicID, r.Format)

			ht := sha256.New()
			ht.Write([]byte(name))
			hash := hex.EncodeToString(ht.Sum(nil))
			reader, err := sd.NewReader(hash, "upload/")
			// close reader if object exists
			if err == nil {
				reader.Close()
			} else {
				// reader is nil, get and save
				resp, err := http.Get(r.URL)
				t2 := time.Now()
				if err != nil {
					log.Println("Failed to download", r.URL)
					return
				}
				body, err := ioutil.ReadAll(resp.Body)
				resp.Body.Close()

				if err != nil {
					log.Println("Failed to parse body", r.URL)
					return
				}
				sd.Write(body, hash, "upload/")
				fmt.Printf("\nget %s, write %s - %s ", t2.Sub(t1), time.Since(t2), r.URL)
			}

		}(r, sd)
	}
	t3 := time.Now()
	wg.Wait()
	fmt.Printf("\nwait %s\n", time.Since(t3))

}

// Import browse all images in userspace and save them in storage
// This function will be running until cloudinary returns empty next_cursor
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
		resp.Body.Close()

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
