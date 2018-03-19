[![Go Report Card](https://goreportcard.com/badge/github.com/yangl900/log2oms)](https://goreportcard.com/report/github.com/yangl900/log2oms) [![Build Status](https://travis-ci.org/yangl900/log2oms.svg?branch=master)](https://travis-ci.org/yangl900/log2oms)

# log2oms
A super tiny agent (binary 5MB, container 12MB) that pushs app logs to Azure Log Analytics (OMS)

# Why we need this
I have been exploring options to push container logs to a remote storage like Log Analytics. A few available options are:

* Use OMS container (https://hub.docker.com/r/microsoft/oms/) to push logs. However, 1) this solution requires running the OMS container as privileged, so it won't work on services like ACI. 2) the size of the image (307MB)...isn't very nice.
* Install OMS agent into my app container. I tried, and realized 1) it comes with lots of dependencies, python etc. 2) it doesn't support alpine. 3) Size of just the installer (omsagent-1.4.4-210.universal.x64.sh 110MB), isn't very container fridenly. Since I only want to upload logs, most of the dependencies are really unnecessary.

So I implemented this tiny agent that uses Log Analytics data collector API (https://docs.microsoft.com/en-us/azure/log-analytics/log-analytics-data-collector-api), and make it container friendly.

# How to use it
## Sidecar

The docker image is published at https://hub.docker.com/r/yangl/log2oms/

The best way to use log2oms is by adopting the "[sidecar](https://docs.microsoft.com/en-us/azure/architecture/patterns/sidecar)" pattern. Having log2oms container (image `yangl/log2oms`) run as a "sidecar" of your app container, and use a shared volume to read the app logs and uplaod to log2oms.

Take a nginx web server as example, you simply run `yangl/log2oms` as another container and shares the nginx /var/log/nginx volume. log2oms will tail the nginx logs and upload to Log Analytics automatically.

```
  +-----------------------------+
  |              |              |
  |    NGINX     |     log2oms  |
  |              |              |
  +-----------------------------+
  |        (shared volume)      |
  |        /var/log/nginx       |
  +-----------------------------+
```

The log2oms container requires only 4 environment variables to run:

* `LOG2OMS_WORKSPACE_ID` This is the workspace ID of Log Analytics.
* `LOG2OMS_WORKSPACE_SECRET` This is the secret of your workspace, you can find it from "Advanced Settings" in Azure portal.
* `LOG2OMS_LOG_FILE` This is the log file to tail and upload. Right now only support 1 file, in nginx case, this will be `access.log`
* `LOG2OMS_LOG_TYPE` This is the table you want logs upload to. Note that LogAnalytics will add a postfix `_CL` to this name. so if we have `nginx` here, in LogAnalytics the table will be `nginx_CL`.

And that's it. No changes needed from app container.

More flags:
* `LOG2OMS_METADATA_*` This is an environment variable prefix for log metadata. The metadata will be sent to Log Analytics for every log message. This is useful if you have multiple replicas sending logs and want to differentiate them. For example, set `LOG2OMS_METADATA_Location=WestUS` and `LOG2OMS_METADATA_Role=Frontend`, logs in Analytics will have 2 more columns `Location` and `Role`.

## Sample for Kubernetes
`samples/kubernetes/deploy.yaml` is a sample yaml how to deploy an nginx server with log2oms as a sidecar. 

To try the sample:
1. Run `kubectl create -f samples/kubernetes/deploy.yaml` to create the deployment. 
2. curl the pod IP on port 80 to generate a few lines of nginx logs.
3. `kubectl logs {pod-name} log2oms` should show some agent logs like following
```
[LOG2OMS][2018-03-17T04:21:36Z] Start tail logs from: /logs/access.log
[2018-03-17T04:22:56Z] 10.244.0.1 - - [17/Mar/2018:04:22:56 +0000] "GET / HTTP/1.1" 200 612 "-" "curl/7.47.0" "-"
[LOG2OMS][2018-03-17T04:23:01Z] Posted 1 messages.
```
4. Wait a few minutes to let LogAnalytics process, then you can query `nginx_access_CL | take 100` in LogAnalytics to see the nginx access logs.

## Sample for Azure Container Instances
1. Click the deploy-to-azure button in [ACI sample page](https://github.com/yangl900/log2oms/tree/master/samples/azure-container-instance) to create an nginx container in ACI with the log2oms sidecar. 
2. After deployment succeed (should be a few seconds), access the container public IP address to generate a few lines of logs. Use Azure portal or [Cloud Shell](https://shell.azure.com) command `az container show -g {resource-group} -n {container-group-name}` to find out the public IP address.
3. Wait a few minutes to let LogAnalytics process, then you can query `nginx_access_CL | take 100` in LogAnalytics to see the nginx access logs.

## Run binary directly
You can also run log2oms as a binary directly, anywhere, either you need upload existing log files or tail a file and keep uploading.

```bash
#! /bin/bash
export LOG2OMS_WORKSPACE_ID={workspace-id}
export LOG2OMS_WORKSPACE_SECRET={workspace-secret}
export LOG2OMS_LOG_FILE={path-to-logfile}
export LOG2OMS_LOG_TYPE={log-table-name}

curl -sL https://github.com/yangl900/log2oms/releases/download/v0.1.0/log2oms_linux_64-bit.tar.gz | tar xz && ./log2oms
```

# Future improvements
* Send a heartbeat signal to log analytics so you know when it is working / stop working.
* Handle SIGTERM to flush out logs before termination.
* Exit on a termination signal file. This will be useful for task containers so the sidecar can stop automatically.
