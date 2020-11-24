package version

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/go-git/go-git/v5/"
)

var (
	// ErrBranchExists an error stating the specified branch already exists
	ErrInvalidMajorVersion = errors.New("major version should be in the \"<int>.x\" format")
)

var validMajorVersion = regexp.MustCompile(`v?\d+\.x`)

const majorVersion = "0.x"

func isMajorVersion(ver string) error {
	if !validMajorVersion.MatchString(ver) {
		return ErrInvalidMajorVersion
	}
	return nil
}

func GetVersion() string {
	dir, err := os.Getwd()

	repo, err := PlainOpen(dir)
	if errors.Is(err, git.ErrRepositoryNotExists) {
		return majorVersion
	}
	repo.Head()
	repo.Describe(repo.Head())

}
