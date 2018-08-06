param (
    [string]$repo = ""
)

$tag = "log2oms:nanoserver-1803"

if ($repo -ne "") {
    $tag = "$repo/$tag"
}

docker build -f .\Dockerfile-Windows -t $tag .
