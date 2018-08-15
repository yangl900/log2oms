param (
    [Parameter(Mandatory=$true)]
    [string]$repo
)

$tag = "log2oms:nanoserver-1803"
docker push "$repo/$tag"