#!/bin/bash
# Update VERSION then run to update the repo (local and remote)
# with the tags. Accepts a quotes string for the commit message.
# For packages with dependencies, the tip of master and the tag will
# be the same code, but the tip will have the updates for packages with
# dependencies.
VERSION="v2.0.5"

git add -A
git commit -m "{$1}"
git push origin

# git tag archiveh/"${VERSION}"    HAS DEPENDENCY
git tag cryptoh/"${VERSION}"
# git tag databaseh/"${VERSION}"   HAS DEPENDENCY
git tag encodingh/"${VERSION}"
git tag logh/"${VERSION}"
git tag mathh/"${VERSION}"
# git tag neth/"${VERSION}"        HAS DEPENDENCY
git tag osh/"${VERSION}"
git tag slicesh/"${VERSION}"
git tag testingh/"${VERSION}"

git push origin --tags

cd archiveh
go get -u github.com/paulfdunn/go-helper/cryptoh/v2@"${VERSION}"
go get -u github.com/paulfdunn/go-helper/testingh/v2@"${VERSION}"
cd ../

cd databaseh
go get -u github.com/paulfdunn/go-helper/osh/v2@"${VERSION}"
cd ../

cd neth
go get -u github.com/paulfdunn/go-helper/osh/v2@"${VERSION}"
cd ../

git add -A
git commit -m 'Update packages with dependencies'
git push origin
git tag archiveh/"${VERSION}"
git tag databaseh/"${VERSION}"
git tag neth/"${VERSION}"
git push origin --tags

echo "\n\nUse the below to update modules"
echo "go get -u github.com/paulfdunn/go-helper/archiveh/v2@${VERSION}; go mod tidy"
echo "go get -u github.com/paulfdunn/go-helper/cryptoh/v2@${VERSION}; go mod tidy"
echo "go get -u github.com/paulfdunn/go-helper/databaseh/v2@${VERSION}; go mod tidy"
echo "go get -u github.com/paulfdunn/go-helper/encodingh/v2@${VERSION}; go mod tidy"
echo "go get -u github.com/paulfdunn/go-helper/logh/v2@${VERSION}; go mod tidy"
echo "go get -u github.com/paulfdunn/go-helper/mathh/v2@${VERSION}; go mod tidy"
echo "go get -u github.com/paulfdunn/go-helper/neth/v2@${VERSION}; go mod tidy"
echo "go get -u github.com/paulfdunn/go-helper/osh/v2@${VERSION}; go mod tidy"
echo "go get -u github.com/paulfdunn/go-helper/slicesh/v2@${VERSION}; go mod tidy"
echo "go get -u github.com/paulfdunn/go-helper/testingh/v2@${VERSION}; go mod tidy"