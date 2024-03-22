#!/bin/bash
VERSION="v1.3.0"
NEXT_VERSION="v1.4.0"

git add -A
git commit $1
git push origin

git tag archiveh/"${VERSION}"
git tag cryptoh/"${VERSION}"
git tag databaseh/"${VERSION}"
git tag encodingh/"${VERSION}"
git tag logh/"${VERSION}"
git tag mathh/"${VERSION}"
git tag neth/"${VERSION}"
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
git tag archiveh/"${NEXT_VERSION}"
git tag archiveh/"${NEXT_VERSION}"
git tag archiveh/"${NEXT_VERSION}"
git push origin --tags

