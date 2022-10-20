# Configmap Watcher
This application provides a kubernetes configmap watcher.
After any changes on files in a directory or a single file in a directory which is mounted throughout a configmap, it will send a SIGKILL to a process which you want.
It is helpful for platforms which need to restart a container after changing the configmap and also if you don't want to use services such as configmap reloader which work by rollout strategy.

## Environment variables
    Name                              | Required        | Description
    ----------------------------------|-----------------|------------
    `FILE_PATH`                       | True            | It should be start with "/", for example "/tmp"
    `PROCESS_NAME`                    | True            | Process name which you want to send SIGKILL

## Requirements
- consider that you should set shareProcessNamespace to true to share processes between containers inside a pod
```
shareProcessNamespace: true
```

## Docker image
You can access to the image from dockerhub:
```
docker push mahdiabbasi/configmap-watcher:v1
```
