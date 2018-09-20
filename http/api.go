package http

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")

// UploadAPIResponse holds the response json model for apiupload
type UploadAPIResponse struct {
	URL   string `json:"url"`
	Error string `json:"error"`
}

func diverseName(name string) (string, error) {
	n := 8
	b := make([]rune, n)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	parts := strings.Split(name, ".")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid name %s", name)
	}
	newName := fmt.Sprintf("%s-%s.%s", parts[0], string(b), parts[1])
	return newName, nil
}

func writeErr(w http.ResponseWriter, msg string, status int) {
	b, err := json.Marshal(&UploadAPIResponse{URL: "", Error: msg})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(status)
		w.Write(b)
	}
}

// APIUpload handles image uploads from API users
func APIUpload(opts *ServerOpts) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if r.Method != "POST" {
			writeErr(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		r.ParseMultipartForm(32 << 20)
		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			log.Printf("can not retrieve file: %v", err)
			writeErr(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		name, err := diverseName(fileHeader.Filename)
		if err != nil {
			log.Printf("can not retrieve file: %v", err)
			writeErr(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		ht := sha256.New()
		ht.Write([]byte(name))
		hash := hex.EncodeToString(ht.Sum(nil))

		opts.StorageDriver.Init()
		reader, err := opts.StorageDriver.NewReader(hash, "upload/")
		if err == nil {
			reader.Close()
			writeErr(w, "Image already exists", http.StatusBadRequest)
			return
		}

		body, err := ioutil.ReadAll(file)
		if err != nil {
			log.Printf("can not retrieve file: %v", err)
			writeErr(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		opts.StorageDriver.Write(body, hash, "upload/")
		finalURL := fmt.Sprintf("https://%s/image/upload/f_auto/%s", opts.Domain, name)
		fmt.Println("Uploaded filename", finalURL)
		b, err := json.Marshal(&UploadAPIResponse{URL: finalURL, Error: ""})
		w.Write(b)
	}
}
