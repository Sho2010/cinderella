package encrypt_test

import (
	"testing"

	"github.com/Sho2010/cinderella/api/v1alpha1"
	"github.com/Sho2010/cinderella/internal/encrypt"
)

func TestFetch(t *testing.T) {
	g := &encrypt.GithubKey{
		Github: v1alpha1.Github{
			User: "sho2010",
		},
	}

	pub, err := g.Fetch()
	if err != nil {
		t.Errorf("expecting bar, got %e", err)
	}
	println(pub.KeyIdString())
}
