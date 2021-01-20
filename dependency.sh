#!/bin/python
cmdexit()
{
if [ $1 -ne "0" ]
then
	exit 1;
fi
}

git config --global url."git@git.code.oa.com:".insteadOf "https://git.code.oa.com/"

if [ $http_proxy="" ];then
git config --global http.proxy http://web-proxy.tencent.com:8080
git config --global https.proxy http://web-proxy.tencent.com:8080

export https_proxy=http://web-proxy.tencent.com:8080
export no_proxy=localhost,127.0.0.1,.oa.com
fi

export GO111MODULE=on


go mod tidy -v
cmdexit $?

rm -r vendor

go mod vendor -v
cmdexit $?


git config --global url."git@git.code.oa.com:".insteadOf "git@git.code.oa.com:"
cmdexit $?
