<p align="center">
<img src="https://i.imgur.com/Wdp5QWg.png" width="750" />
</p>
K8s admission controller written in Golang with Fiber framework. 
</br>
This admission controller target is to ensure k8s best practices are kept.

# Why use Fenrir ?
## Fast
* It runs with Fiber, which is the fastest framework out there.
* It Logs with Zap, which is the fastest logger out there.

## Configurable
the admission controller is configurable:
### Environment vars config:
* `LOG_LVL, default is info.`
* `PORT, default is 8080.`
* `OUTPUT, default is stdout.`
* `CONFIG_POLICY_PATH, default is ./conf.json`

### Policy configuration:
Its based on json file in location - `CONFIG_POLICY_PATH`
</br>
The policy updates at real time, after you change json file.

**json policy sample:**
```json
{
  "pod": {
    "policy_enforcement": true,
    "default_policy": {
      "readiness_liveness": true,
      "default_ns": true,
      "latest_image_tag": false,
      "run_as_non_root": false
    },
    "custom_policies": {}
  },
  "service": {
    "policy_enforcement": true,
    "default_policy": {
      "load_balancer": true,
      "default_ns": true
    },
    "custom_policies": {}
  }
}
```

### Policy fields and validations
**under pod we have:**
* **readiness_liveness -** checks if your pod has liveness & readiness.
* **default_ns -** checks that you dont try to deploy pods on default ns.
* **latest_image_tag -** checks that you dont try to deploy latest image tag.
* **run_as_non_root -** checks that you dont try to run as root.
* **resources -** checks that you state your resource usage.


**uder service we have:**
</br>
* **load_balancer -** checks if service is of type LoadBalancer.
* **default_ns -** checks that you dont try to deploy pods on default ns.

### Note: you can set different policy for each ns

## Light
* It's written in Golang.

## You can run it anywhere
* you can compile it to statically linked executable, for any OS.

# Features
- [x] Support graceful shutdown.
- [x] Support policy file.
- [x] Support live changes in policy. 
- [x] Pod policy impl.
- [x] Service policy impl.
- [ ] Ingress policy impl.
- [ ] Deployment policy impl.
- [ ] DeploymentConfig policy impl.
- [ ] Route policy impl.