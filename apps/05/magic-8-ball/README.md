## magic-8-ball template

To deploy this application to a WASM-enabled Kubernetes cluster, run

```
npm install

spin build

spin k8s scaffold <image repo>

spin k8s build

spin k8s push

spin k8s deploy
```
