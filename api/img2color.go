package handler

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/joho/godotenv"
	"image"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Response struct {
	RGB string `json:"RGB"`
}

var ctx = context.Background()
var kvEnable bool
var kvURL string
var kvToken string

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	kvEnable, _ = strconv.ParseBool(os.Getenv("KV_ENABLE"))
	if kvEnable {
		kvURL = os.Getenv("KV_REST_API_URL")
		kvToken = os.Getenv("KV_REST_API_TOKEN")
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if !checkReferer(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	imgURL := r.URL.Query().Get("img")
	if imgURL == "" {
		http.Error(w, "img parameter is required", http.StatusBadRequest)
		return
	}

	var color string
	var err error
	if kvEnable {
		color, err = getColorFromKV(imgURL)
		if err != nil {
			color, err = getColorFromImageURL(imgURL)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			setColorToKV(imgURL, color)
		}
	} else {
		color, err = getColorFromImageURL(imgURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	res := Response{RGB: color}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func checkReferer(r *http.Request) bool {
	referer := r.Header.Get("Referer")
	allowedReferers := strings.Split(os.Getenv("ALLOWED_REFERERS"), ",")
	for _, allowedReferer := range allowedReferers {
		if allowedReferer == "*" || strings.HasSuffix(referer, allowedReferer) {
			return true
		}
	}
	return false
}

func getColorFromKV(imgURL string) (string, error) {
	hasher := sha256.New()
	hasher.Write([]byte(imgURL))
	key := hex.EncodeToString(hasher.Sum(nil))
	req, err := http.NewRequest("GET", kvURL+"/"+key, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+kvToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", fmt.Errorf("Not Found")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func setColorToKV(imgURL string, color string) {
	hasher := sha256.New()
	hasher.Write([]byte(imgURL))
	key := hex.EncodeToString(hasher.Sum(nil))
	req, err := http.NewRequest("PUT", kvURL+"/"+key, bytes.NewBuffer([]byte(color)))
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+kvToken)
	req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
}

func getColorFromImageURL(imgURL string) (string, error) {
	resp, err := http.Get(imgURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return "", err
	}

	img = imaging.Resize(img, 1, 1, imaging.Box)
	dc := gg.NewContext(1, 1)
	dc.DrawImage(img, 0, 0)
	rVal, g, b, _ := dc.Image().At(0, 0).RGBA()
	color := fmt.Sprintf("#%02X%02X%02X", uint8(rVal>>8), uint8(g>>8), uint8(b>>8))

	return color, nil
}
