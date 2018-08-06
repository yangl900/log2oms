param (
    [Parameter(Mandatory=$true)]
    [string]$repo
)

docker push "$repo/log2oms:nanoserver-1803"