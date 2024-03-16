find  . -type d -not -path "./.git*" | xargs -I % -L1 sh -c 'cd %; pwd; go get -u; cd ..'
find  . -type d -not -path "./.git*" | xargs -I % -L1 sh -c 'cd %; pwd; go mod tidy; cd ..'