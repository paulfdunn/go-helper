#!/bin/bash
# Update VERSION then run to update the repo (local and remote)
# with the tags. Accepts a quotes string for the commit message.
# For packages with dependencies, the tip of master and the tag will
# be the same code, but the tip will have the updates for packages with
# dependencies.
VERSION="v2.0.0"

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
go get -u github.com/paulfdunn/go-helper/cryptoh@"${VERSION}"
go get -u github.com/paulfdunn/go-helper/testingh@"${VERSION}"
cd ../

cd databaseh
go get -u github.com/paulfdunn/go-helper/osh@"${VERSION}"
cd ../

cd neth
go get -u github.com/paulfdunn/go-helper/osh@"${VERSION}"
cd ../

git add -A
git commit -m 'Update packages with dependencies'
git push origin
git tag archiveh/"${VERSION}"
git tag databaseh/"${VERSION}"
git tag neth/"${VERSION}"
git push origin --tags

echo "\n\nUse the below to update modules"
echo "go get -u github.com/paulfdunn/go-helper/archiveh@${VERSION}; go mod tidy"
echo "go get -u github.com/paulfdunn/go-helper/cryptoh@${VERSION}; go mod tidy"
echo "go get -u github.com/paulfdunn/go-helper/databaseh@${VERSION}; go mod tidy"
echo "go get -u github.com/paulfdunn/go-helper/encodingh@${VERSION}; go mod tidy"
echo "go get -u github.com/paulfdunn/go-helper/logh@${VERSION}; go mod tidy"
echo "go get -u github.com/paulfdunn/go-helper/mathh@${VERSION}; go mod tidy"
echo "go get -u github.com/paulfdunn/go-helper/neth@${VERSION}; go mod tidy"
echo "go get -u github.com/paulfdunn/go-helper/osh@${VERSION}; go mod tidy"
echo "go get -u github.com/paulfdunn/go-helper/slicesh@${VERSION}; go mod tidy"
echo "go get -u github.com/paulfdunn/go-helper/testingh@${VERSION}; go mod tidy"