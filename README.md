<p align="center">
<img src="https://i.imgur.com/Wdp5QWg.png" width="600" />
</p>
<h1 align="center">Fenrir</h1>

K8s admission controller written in Golang with Fiber framework. 
</br>
This admission controller target is to ensure k8s best practices are kept.

# Why use Fenrir ?
## Fast
* It runs with Fiber, which is the fastest framework out there.
* It Logs with Zap, which is the fastest logger out there.

## Configurable
the admission controller is configurable:
### General config:
* `LOG_LVL, default is info.`
* `PORT, default is 8080.`
* `OUTPUT, default is stdout.`
* `LOGGER_NAME, default is github.com.matankila.fenrir.logger`
### Pod Policy config:
* `POD_LIVENESS_READINESS_CHECK, default is true.`
* `POD_RESTRICTED_NS_CHECK, default is true.`
* `POD_RUN_AS_NON_ROOT_CHECK, default is false.`
* `POD_LATEST_IMAGE_TAG_CHECK, default is false.`

## Light
* It's written in Golang.

## You can run it anywhere
* you can compile it to statically linked executable, for any OS. 

# State
- [x] Pod policy impl.
- [ ] Service policy impl.
- [ ] Ingress policy impl.
- [ ] Deployment policy impl.
- [ ] DeploymentConfig policy impl.
- [ ] Route policy impl.