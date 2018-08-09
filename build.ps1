param (
    [string]$repo = "myplooploops"
)

docker build -f .\windows-container\windows-build\windows-go\Dockerfile -t "$repo/golang:windowsservercore-1803" .
docker build -f .\windows-container\Dockerfile-1709 -t "$repo/log2oms:nanoserver-1709" .
docker build -f .\windows-container\Dockerfile-1803 -t "$repo/log2oms:nanoserver-1803" .
