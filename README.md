* assumes that most usage is in a git repository
* falls back to a sensible constant
* deduplicates version tagging by relying on tag precedence of:
    1. `git tag`
    1. `HEAD` sha
    1. hardcoded major version constant

```
major_version = 1.x
version = major_version


if not inGitRepository():
    return version

git_tag = git describe --tags --exact-match HEAD

if git_tag is set:
    version = git_tag
else:
    version = git rev-parse HEAD

if not matchesMajorVersion(version, major_version):
    exit 1

return version
```
