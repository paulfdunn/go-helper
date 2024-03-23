#!/bin/bash
# Update VERSION and VERSION_WITH_DEPENDENCIES, then run to update update the repo (local and remote)
# with the tags. Accepts a quotes string for the commit message
# Note that the versioning scheme is that with no interdependencies get the odd minor version, then
# packages with interdependencies get the next higher even version.
VERSION="v1.3.7"
VERSION_WITH_DEPENDENCIES="v1.4.7"

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
git tag archiveh/"${VERSION_WITH_DEPENDENCIES}"
git tag databaseh/"${VERSION_WITH_DEPENDENCIES}"
git tag neth/"${VERSION_WITH_DEPENDENCIES}"
git push origin --tags

