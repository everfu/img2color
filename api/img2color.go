// 如果在本地运行请改为 main
package handler

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"image"
	"image/draw"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Response struct {
	RGB string `json:"RGB"`
}

var ctx = context.Background()
var kvEnable bool
var kvURL string
var kvToken string

var rdb *redis.Client

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	kvEnable, _ = strconv.ParseBool(os.Getenv("KV_ENABLE"))
	if kvEnable {
		rdb = redis.NewClient(&redis.Options{
			Addr:     kvURL,
			Password: kvToken,
			DB:       0,
		})
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

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
		allowedReferer = strings.ReplaceAll(allowedReferer, ".", "\\.")
		allowedReferer = strings.ReplaceAll(allowedReferer, "*", ".*")
		match, _ := regexp.MatchString(allowedReferer, referer)
		if match {
			return true
		}
	}
	return false
}

func getColorFromKV(imgURL string) (string, error) {
	hasher := sha256.New()
	hasher.Write([]byte(imgURL))
	key := hex.EncodeToString(hasher.Sum(nil))

	color, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("Not Found")
	} else if err != nil {
		return "", err
	}

	return color, nil
}

func setColorToKV(imgURL string, color string) {
	hasher := sha256.New()
	hasher.Write([]byte(imgURL))
	key := hex.EncodeToString(hasher.Sum(nil))

	err := rdb.Set(ctx, key, color, 0).Err()
	if err != nil {
		fmt.Println(err)
		return
	}
}

var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func getColorFromImageURL(imgURL string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, imgURL, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	buf := bufferPool.Get().(*bytes.Buffer)
	defer bufferPool.Put(buf)
	buf.Reset()
	buf.Write(bodyBytes)

	var img image.Image
	img, _, err = image.Decode(buf)
	if err != nil {
		img, err = webp.Decode(buf)
		if err != nil {
			return "", err
		}
	}

	img = imaging.Resize(img, 1, 1, imaging.Box)
	dc := image.NewRGBA(image.Rect(0, 0, 1, 1))
	draw.Draw(dc, dc.Bounds(), img, img.Bounds().Min, draw.Src)
	rVal, g, b, _ := dc.At(0, 0).RGBA()
	color := fmt.Sprintf("#%02X%02X%02X", uint8(rVal>>8), uint8(g>>8), uint8(b>>8))

	return color, nil
}
func main() {
	http.HandleFunc("/", Handler)
	http.ListenAndServe(":8080", nil)
}
