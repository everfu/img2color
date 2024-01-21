package handler

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"image"
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

var rdb *redis.Client

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	kvEnable, _ = strconv.ParseBool(os.Getenv("KV_ENABLE"))
	if kvEnable {
		rdb = redis.NewClient(&redis.Options{
			Addr:     kvURL,   // Redis server URL
			Password: kvToken, // no password set
			DB:       0,       // use default DB
		})
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
