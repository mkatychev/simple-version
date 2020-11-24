package version

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

// Repo struct wrapps Repository class from go-git to add a tag map used to perform queries when describing.
type Repo struct {
	TagsMap map[plumbing.Hash]*plumbing.Reference
	*git.Repository
}

// PlainOpen opens a git repository from the given path. It detects if the
// repository is bare or a normal one. If the path doesn't contain a valid
// repository ErrRepositoryNotExists is returned
func PlainOpen(path string) (*Repo, error) {
	r, err := git.PlainOpen(path)
	return &Repo{
		make(map[plumbing.Hash]*plumbing.Reference),
		r,
	}, err
}

func (g *Repo) getTagMap() error {
	tags, err := g.Tags()
	if err != nil {
		return err
	}

	err = tags.ForEach(func(t *plumbing.Reference) error {
		g.TagsMap[t.Hash()] = t
		return nil
	})

	return err
}

// Describe the reference as 'git describe --tags' will do
func (g *Repo) Describe(reference *plumbing.Reference) (string, error) {
	// Fetch the reference log
	commits, err := g.Log(&git.LogOptions{
		From:  reference.Hash(),
		Order: git.LogOrderCommitterTime,
	})

	// Build the tag map
	err = g.getTagMap()
	if err != nil {
		return "", err
	}

	// Search the tag
	var tag *plumbing.Reference
	var count int
	err = commits.ForEach(func(c *object.Commit) error {
		if t, ok := g.TagsMap[c.Hash]; ok {
			tag = t
		}
		if tag != nil {
			return storer.ErrStop
		}
		count++
		return nil
	})
	if count == 0 {
		return fmt.Sprint(tag.Name().Short()), nil
	} else {
		return fmt.Sprintf("%v-%v-%v",
			tag.Name().Short(),
			count,
			tag.Hash().String()[0:8],
		), nil
	}
}

// func main() {
//     dir, err := os.Getwd()

//     repo, err := git.PlainOpen(dir)

//     revision := "origin/master"

//     revHash, err := repo.ResolveRevision(plumbing.Revision(revision))
//     CheckIfError(err)
//     revCommit, err := repo.CommitObject(*revHash)

//     CheckIfError(err)

//     headRef, err := repo.Head()
//     CheckIfError(err)
//     // ... retrieving the commit object
//     headCommit, err := repo.CommitObject(headRef.Hash())
//     CheckIfError(err)

//     isAncestor, err := headCommit.IsAncestor(revCommit)

//     CheckIfError(err)

//     fmt.Printf("Is the HEAD an IsAncestor of origin/master? : %v\n", isAncestor)
// }
