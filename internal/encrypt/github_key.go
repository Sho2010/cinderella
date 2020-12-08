package encrypt

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Sho2010/cinderella/api/v1alpha1"
	"golang.org/x/crypto/openpgp/packet"
)

type GithubKey struct {
	v1alpha1.Github
	publicKey *packet.PublicKey
	lastFetch time.Time
}

var (
	CacheExpirationSecond = 300
)

func (g *GithubKey) url() string {
	// TODO: Github enterprise
	return fmt.Sprintf("https://github.com/%s.gpg", g.User)
}

// Fetch is implements PublicKeySource interface
// fetch key from 'github.com/<user>.gpg' URL and store cache it
func (g *GithubKey) Fetch() (*packet.PublicKey, error) {
	c := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := c.Get(g.url())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// TODO: check status code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf(fmt.Sprintf("HTTP STATUS : %d", resp.StatusCode))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	key, err := g.parseGithubGpg(string(body))
	if err != nil {
		return nil, err
	}

	pub, err := decodePublicKey(key)
	g.publicKey = pub
	g.lastFetch = time.Now()

	return g.publicKey, nil
}

func (g *GithubKey) cahceValid() bool {
	d := time.Duration(CacheExpirationSecond) * time.Second
	return g.lastFetch.Add(d).Before(time.Now())
}

func (g *GithubKey) parseGithubGpg(keyStr string) (string, error) {
	//implements me
	return keyStr, nil
}
