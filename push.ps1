#MSMQ
param (
    [string]$repo = "myplooploops"
)

docker push "$repo/golang:windowsservercore-1803"
docker push "$repo/log2oms:nanoserver-1709"
docker push "$repo/log2oms:nanoserver-1803"
