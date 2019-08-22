package weather

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

// UoM describes imperial and metric systems..
type UoM int

const (
	Imperial UoM = iota
	Metric
)

const (
	RssUrl = "https://weather-ydn-yql.media.yahoo.com/forecastrss"
)

var (
	// MinUpdateTimeoutSeconds defines number of seconds before next actual update.
	MinUpdateTimeoutSeconds = int64(5 * 60)
)

// YahooWeatherProvider defines Yahoo Weather wrapper.
type YahooWeatherProvider struct {
	appID        string
	clientID     string
	clientSecret string
	compositeKey string

	lastLocation     string
	lastLocationNorm string
	lastUnit         UoM
	lastUnitStr      string
	lastQueryTime    int64

	lastData *Weather
}

// NewProvider constructs a new Yahoo provider.
func NewProvider(appID string, clientID string, clientSecret string) *YahooWeatherProvider {
	yw := &YahooWeatherProvider{
		appID:        appID,
		clientID:     clientID,
		clientSecret: clientSecret,
		compositeKey: url.QueryEscape(clientSecret) + "&",
		lastUnitStr:  "f",
		lastUnit:     Imperial,

		lastData: &Weather{},
	}
	return yw
}

// Query gets current weather at the specified location.
// If location is the same as in previous request and now() - last_request_time() <  MinUpdateTimeoutSeconds,
// previous result used.
func (provider *YahooWeatherProvider) Query(location string, unit UoM) (*Weather, error) {
	var err error = nil

	if location != provider.lastLocation || unit != provider.lastUnit ||
		time.Now().UTC().Unix()-provider.lastQueryTime > MinUpdateTimeoutSeconds {
		provider.lastLocation = location
		provider.lastLocationNorm = strings.ReplaceAll(strings.ToLower(location), ", ", ",")
		provider.lastUnit = unit
		if unit == Metric {
			provider.lastUnitStr = "c"
		} else {
			provider.lastUnitStr = "f"
		}

		err = provider.update()
	}

	return provider.lastData, err
}

// Performs actual update.
func (provider *YahooWeatherProvider) update() error {
	sign, err := provider.getAuth()
	if err != nil {
		return err
	}

	url_ := fmt.Sprintf("%s?location=%s&u=%s&format=json", RssUrl,
		provider.lastLocationNorm, provider.lastUnitStr)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url_, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Yahoo-App-Id", provider.appID)
	req.Header.Set("Authorization", sign)

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("wrong HTTP status: %d", res.StatusCode))
	}

	dec := json.NewDecoder(res.Body)
	weather := &Weather{}
	err = dec.Decode(weather)
	if err != nil {
		provider.lastData = nil
		return err
	}

	provider.lastData = weather
	provider.lastQueryTime = time.Now().UTC().Unix()

	return nil
}

// Generating OAuth1 signature.
func (provider *YahooWeatherProvider) getAuth() (string, error) {
	nonce, err := getNonce()
	if err != nil {
		return "", err
	}
	oauth := map[string]string{
		"oauth_consumer_key":     provider.clientID,
		"oauth_nonce":            nonce,
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp":        getTimestamp(),
		"oauth_version":          "1.0",
	}

	query := map[string]string{
		"location": provider.lastLocationNorm,
		"format":   "json",
		"u":        provider.lastUnitStr,
	}

	merged := make(map[string]string)

	for k, v := range oauth {
		merged[k] = v
	}

	for k, v := range query {
		merged[k] = v
	}

	keys := make([]string, 0)

	for k, _ := range merged {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	sortedParams := make([]string, 0)

	for _, v := range keys {
		sortedParams = append(sortedParams, fmt.Sprintf("%s=%s", url.QueryEscape(v), url.QueryEscape(merged[v])))
	}

	baseString := fmt.Sprintf("GET&%s&%s", url.QueryEscape(RssUrl),
		url.QueryEscape(strings.Join(sortedParams, "&")))

	h := hmac.New(sha1.New, []byte(provider.clientSecret+"&"))
	_, err = h.Write([]byte(baseString))
	if err != nil {
		return "", err
	}

	oauth["oauth_signature"] = base64.StdEncoding.EncodeToString(h.Sum(nil))

	headerParams := make([]string, 0)
	for k, v := range oauth {
		headerParams = append(headerParams, fmt.Sprintf("%s=\"%s\"", k, v))
	}

	return fmt.Sprintf("OAuth %s", strings.Join(headerParams, ", ")), nil
}

// Generating 32-byte nonce.
func getNonce() (string, error) {
	h := md5.New()
	_, err := h.Write(h.Sum([]byte(time.Now().Format(time.RFC3339Nano))))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// Obtaining current timestamp in UTC.
func getTimestamp() string {
	return strconv.FormatInt(time.Now().UTC().Unix(), 10)
}
