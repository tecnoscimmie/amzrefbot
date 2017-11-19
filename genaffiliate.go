package main

import (
	"errors"
	"math/rand"
	"net/url"
	"strings"
	"time"
)

// GenAffiliate generates an affiliate link from a standard Amazon link
func (r Refs) GenAffiliate(inputURL string) (string, string, error) {
	u, err := url.Parse(inputURL)
	if err != nil {
		return "", "", err
	}

	if u.Host != "www.amazon.it" {
		return "", "", errors.New("this is not an amazon.it link")
	}

	randomInt := getRandomIntMax(len(r.ReferralCodes))
	randomCode := r.ReferralCodes[randomInt].Code
	randomUser := r.ReferralCodes[randomInt].AssociatedUser

	newPath, err := generateNewPath(u.Path)
	if err != nil {
		return "", "", errors.New("cannot generate path url for " + inputURL + " because: " + err.Error())
	}
	u.Path = newPath
	q := url.Values{}
	q.Set("tag", randomCode)
	u.RawQuery = q.Encode()

	return u.String(), randomUser, nil
}

// generateNewPath generates a new URL path from the given Amazon.it URL.
// Amazon.it URLs contain the product SKU in the URL path itself, so we
// iterate through it, and as soon as we find the "product" or "dp" strings,
// we know that the next array item is the SKU.
func generateNewPath(path string) (string, error) {
	// split up the path by "/"
	pathArray := strings.Split(path, "/")

	var sku string
	// iterate until we find "product" or "dp", or item length is 10 (ASINs)
	for _, item := range pathArray {
		if len(item) == 10 {
			sku = item
			break
		}
	}

	// if we didn't find neither "product" or "dp", error out
	if sku == "" {
		return "", errors.New("product sku not found")
	}

	// return the new path
	return "/dp/" + sku, nil
}

func getRandomIntMax(max int) int {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	return random.Intn(max)
}
