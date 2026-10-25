package onepassword

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	maxRateLimitRetries = 5
	maxRetryWait        = 60 * time.Second
)

// apiTransport sends requests to the 1Password Events and Users APIs, waiting out 429s.
type apiTransport struct {
	httpClient *http.Client
	sleep      func(context.Context, time.Duration) error
	now        func() time.Time
}

func newAPITransport() apiTransport {
	return apiTransport{
		httpClient: &http.Client{Timeout: 60 * time.Second},
		sleep:      sleepContext,
		now:        time.Now,
	}
}

// do builds a fresh request per attempt because a sent request body cannot be replayed.
func (t apiTransport) do(ctx context.Context, newRequest func() (*http.Request, error)) (*http.Response, error) {
	for attempt := 0; ; attempt++ {
		req, err := newRequest()
		if err != nil {
			return nil, err
		}
		resp, err := t.httpClient.Do(req.WithContext(ctx))
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusTooManyRequests || attempt >= maxRateLimitRetries {
			return resp, nil
		}
		wait := retryAfter(resp.Header.Get("Retry-After"), t.now(), attempt)
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
		if err := t.sleep(ctx, wait); err != nil {
			return nil, err
		}
	}
}

// retryAfter reads Retry-After as seconds or an HTTP date, and falls back to exponential backoff.
func retryAfter(header string, now time.Time, attempt int) time.Duration {
	wait := time.Duration(1<<attempt) * time.Second
	header = strings.TrimSpace(header)
	if secs, err := strconv.Atoi(header); err == nil && secs >= 0 {
		wait = time.Duration(secs) * time.Second
	} else if at, err := http.ParseTime(header); err == nil {
		wait = max(at.Sub(now), 0)
	}
	return min(wait, maxRetryWait)
}

func sleepContext(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

type apiError struct {
	API        string
	StatusCode int
	Message    string
}

func (e *apiError) Error() string {
	if e.Message == "" {
		return fmt.Sprintf("1Password %s API returned HTTP %d", e.API, e.StatusCode)
	}
	return fmt.Sprintf("1Password %s API returned HTTP %d: %s", e.API, e.StatusCode, e.Message)
}

func newAPIError(api string, resp *http.Response) *apiError {
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	return &apiError{API: api, StatusCode: resp.StatusCode, Message: errorMessage(body)}
}

// errorMessage understands the error bodies of both APIs: Events uses {"Error":{"Message"}} or
// {"status","message"}, and the Users API and its token endpoint use {"error","error_description"}.
func errorMessage(body []byte) string {
	var parsed struct {
		Error            json.RawMessage `json:"error"`
		ErrorDescription string          `json:"error_description"`
		Message          string          `json:"message"`
		UpperError       struct {
			Message string `json:"Message"`
		} `json:"Error"`
	}
	if json.Unmarshal(body, &parsed) != nil {
		return strings.TrimSpace(string(body))
	}
	var errorString string
	_ = json.Unmarshal(parsed.Error, &errorString)
	var errorObject struct {
		Message string `json:"message"`
	}
	_ = json.Unmarshal(parsed.Error, &errorObject)

	for _, candidate := range []string{parsed.ErrorDescription, parsed.Message, parsed.UpperError.Message, errorObject.Message, errorString} {
		if candidate != "" {
			return candidate
		}
	}
	return strings.TrimSpace(string(body))
}
