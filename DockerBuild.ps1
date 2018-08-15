param (
    [string]$repo = ""
)

$tag = "log2oms:nanoserver-1803"
$path = ".\Dockerfile-Windows-1803"

if ($repo -ne "") {
    $tag = "$repo/$tag"
}

docker build -f $path -t $tag .
