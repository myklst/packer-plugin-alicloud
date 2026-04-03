package common

import (
	"math/rand"
	"net/url"
	"strings"
	"time"

	aliclient "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	alitea "github.com/alibabacloud-go/tea/tea"
)

func RandomString(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}
	return string(b)
}

func IsRetryableError(err error) (bool, error) {
	// Transient condition while CloudAssistant is initializing.
	if strings.Contains(err.Error(), "CloudAssistant.NotReady") || strings.Contains(err.Error(), "CloudAssistant not ready") {
		return true, err
	}
	if aliErr, ok := err.(*aliclient.AlibabaCloudError); ok {
		return false, aliErr
	}
	if aliErr, ok := err.(*aliclient.ServerError); ok {
		return false, aliErr
	}
	if aliErr, ok := err.(*aliclient.ClientError); ok {
		return false, aliErr
	}
	if aliErr, ok := err.(*alitea.SDKError); ok {
		return false, aliErr
	}
	if urlErr, ok := err.(*url.Error); ok {
		if !urlErr.Timeout() && !urlErr.Temporary() {
			return false, urlErr
		}
	}
	return true, err
}

func NilOrString(s string) *string {
	if s == "" {
		return nil
	}
	return alitea.String(s)
}

func NilOrStringSlice(s ...string) []*string {
	if len(s) == 0 {
		return nil
	}
	return alitea.StringSlice(s)
}

func NilOrBool(b bool) *bool {
	if !b {
		return nil
	}
	return alitea.Bool(b)
}

func NilOrInt64(i int64) *int64 {
	if i == 0 {
		return nil
	}
	return alitea.Int64(i)
}
