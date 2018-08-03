#MSMQ
param (
    [string]$repo = "myrepo"
)

docker push "$repo/golang:windowsservercore-1803"
docker push "$repo/log2oms:nanoserver-1803"