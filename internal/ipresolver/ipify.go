package ipresolver

import (
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

const (
	IPV4_API_URL = "https://api.ipify.org"
	IPV6_API_URL = "https://api64.ipify.org"
)

func GetIPv4() (string, error) {
	return getIp(IPV4_API_URL)
}

func GetIPv6() (string, error) {
	return getIp(IPV6_API_URL)
}

func getIp(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		return "", errors.New(fmt.Sprintf("request to %s returned status code %d", url, res.StatusCode))
	}

	defer res.Body.Close()
	ip, err := io.ReadAll(res.Body)
	if err != nil {
		return "", errors.Wrap(err, "failed to read response body")
	}

	return string(ip), nil
}
