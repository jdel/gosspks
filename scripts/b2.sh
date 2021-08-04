#!/bin/bash -e
pyenv global 3.7
sudo pip3 install --upgrade pip
sudo pip3 install --ignore-installed b2
b2 authorize-account $B2_APPLICATION_KEY_ID $B2_APPLICATION_KEY
for binary in out/*; do
    echo "b2 upload-file --info commit=${TRAVIS_COMMIT:0:8} jdel-builds ${binary} ${TRAVIS_REPO_SLUG#jdel/}/${TRAVIS_BRANCH}/${binary#out/}"
    b2 upload-file --info commit=${TRAVIS_COMMIT:0:8} jdel-builds ${binary} ${TRAVIS_REPO_SLUG#jdel/}/${TRAVIS_BRANCH}/${binary#out/}
done
b2 clear-account