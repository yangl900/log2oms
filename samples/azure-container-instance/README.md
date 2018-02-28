# Upload Azure Container Instances logs to OMS

This template demonstrates a simple use case for uploading nginx access logs to Azure Log Analytics. The nginx container runs in [Azure Container Instances](https://docs.microsoft.com/en-us/azure/container-instances/). Clicking the deploy button will create a nginx container with public IP address and a log2oms container as sidecar that uploads the logs.

<a href="https://portal.azure.com/#create/Microsoft.Template/uri/https%3A%2F%2Fraw.githubusercontent.com%2Fyangl900%2Fmaster%2Fsamples%2Fazure-container-instance%2Fazuredeploy.json" target="_blank">
    <img src="http://azuredeploy.net/deploybutton.png"/>
</a>