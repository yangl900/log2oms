param (
    [string]$repo = "myrepo"
)

docker build -f .\windows-container\windows-build\windows-go\Dockerfile -t "$repo/golang:windowsservercore-1803" .
docker build -f .\windows-container\Dockerfile -t "$repo/log2oms:nanoserver-1803" .
