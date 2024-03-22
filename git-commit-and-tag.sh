VERSION="v1.1.0"
NEXT_VERSION="v1.2.0"

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
git tag archiveh/"${NEXT_VERSION}"
cd ../

cd databaseh
go get -u github.com/paulfdunn/go-helper/osh@"${VERSION}"
git tag archiveh/"${NEXT_VERSION}"
cd ../

cd neth
go get -u github.com/paulfdunn/go-helper/osh@"${VERSION}"
git tag archiveh/"${NEXT_VERSION}"
cd ../

git add -A
git commit $1
git push origin
git push origin --tags

