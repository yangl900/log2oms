# log2oms
A super tiny agent (binary 7MB, container 12MB) that pushs app logs to Azure Log Analytics (OMS)

# Why we need this
I have been exploring options to push container logs to a remote storage like Log Analytics. A few available options are:

* Use OMS container (https://hub.docker.com/r/microsoft/oms/) to push logs. However, 1) this solution requires running the OMS container as privileged. 2) the size of the image (307MB)...isn't very nice.
* Install OMS agent into my app container. I tried, and realized 1) it comes with lots of dependencies, python etc. 2) it doesn't support alpine. 2) Size of just the installer (omsagent-1.4.4-210.universal.x64.sh 110MB), isn't very container fridenly..

Given I simply want logs uploaded, I decided to implement a tiny agent that uses Log Analytics data collector API (https://docs.microsoft.com/en-us/azure/log-analytics/log-analytics-data-collector-api), and keept it container friendly.

# How to use it
## Upload container logs
The best way to use this agent is by adopting the "sidecar" pattern. Having this agent container (yangl/log2oms) run as a "sidecar" of your app container, and use a shared volume to share the app logs to log2oms.

Take a nginx web server as example, you simply run yangl/log2oms as another container and shares the nginx /var/log/nginx volume. Log2OMS will tail the nginx logs and upload to Log Analytics.

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
* `LOG2OMS_WORKSPACE_SECRET` This is the secret of your workspace, you can get from "Advanced Settings" in Azure portal.
* `LOG2OMS_LOG_FILE` This is the log file to tail and upload. Right now only support 1 file, in nginx case, this will be `access.log`
* `LOG2OMS_LOG_TYPE` This is the table you want logs upload to. Note that LogAnalytics will add a postfix `_CL` to this name. so if we have `nginx` here, in LogAnalytics the table will be `nginx_CL`.

And that's it. No changes needed from app container.

