# Telepresence Demo

## Prerequisites

### Install Docker

Install the Docker command line tools (to avoid using Docker desktop):

```
brew install docker
```

Setup the terminal for using the Docker server bundled with Minikube (inside a Linux virtual machine):

```
eval $(minikube -p minikube docker-env)
```

### Install Minikube (with Hyperkit)

```
brew install minikube
brew install hyperkit
```

Create a kubernetes cluster:

```
minikube config set driver hyperkit
minikube start
```

### Install and Connect Telepresence

Install the command line tools for telepresence:

```
brew install datawire/blackbird/telepresence
```

Telepresence requires some in-cluster help to manage and redirect traffic. This helper is called `traffic-manager` and can be installed with Helm.

```
helm repo add datawire  https://app.getambassador.io
helm repo update

kubectl create namespace ambassador
helm install traffic-manager --namespace ambassador datawire/telepresence
```

The telepresence command line must initiate a connection to the cluster.

```
telepresence connect
```

See also https://www.getambassador.io/docs/telepresence/latest/install/helm.

## Use Case: Intercepting a Service

### Background

This repo contains three services which are chained together in this flow:

```
svctop => svcmid => svcbot
```

The `svctop` service has an API endpoint which consumes an endpoint on the `svcmid` service. The `svcmid` service endpoint implementation consumes an endpoing on the `svcbot` service.

When a request is made to the `svctop` endpoint, ultimately the response from `svcbot` is bubbled up through the call chain and returned to the caller of the `svctop` service.


### Running the Bottom Service

```
cd svcbot
make deploy
```

Use the telepresence DNS to test the service. Use the ping endpoint:

```
curl http://svcbot.default:8080/ping
{"message":"pong"}
```

Use the `call` endpoint to trigger the handler:

```
curl http://svcbot.default:8080/call
{"message":"Hi, I am svcbot!"}
```

### Running the Middle Service

```
cd svcmid
make deploy
```

Use the telepresence DNS to test the service. Use the ping endpoint:

```
curl http://svcmid.default:8080/ping
{"message":"pong"}
```

Use the `call` endpoint to trigger the handler:

```
curl http://svcmid.default:8080/call
{"message":"Hi, I am svcbot!"}
```

Notice that the response to the `call` endpoint is the same as the `call` endpoint on `svcbot`. That's because `svcmid` uses the `svcbot` endpoint:

```
r.GET("/call", handler)
```

Where `handler` is making an HTTP call to `svcbot`:

```
resp, err := http.Get("http://svcbot.default:8080/call")
```

### Running the Top Service

```
cd svctop
make deploy
```

Use the telepresence DNS to test the service. Use the ping endpoint:

```
curl http://svctop.default:8080/ping
{"message":"pong"}
```

Use the `call` endpoint to trigger the handler:

```
curl http://svctop.default:8080/call
{"message":"Hi, I am svcbot!"}
```

Notice that the response to the `call` endpoint is the same as the `call` endpoint on `svcbot`. That's because `svctop` uses the `svcmid` endpoint (which in turn uses the `svcbot` endpoint):

```
r.GET("/call", handler)
```

Where `handler` is making an in-cluster HTTP call to `svcmid`:

```
resp, err := http.Get("http://svcmid:8080/call")
```

### Setup the Intercept

Now, intercept the `svcmid` service in the cluster so that traffic is routed to the service running in development on a local laptop.

Get the name of the service's port from kubernetes (e.g. `http`):

```
kubectl get svc svcmid -o yaml | grep -A4 ports:
```

Use that port name in the `telepresence intercept` command to map port `8080` on the local development laptop to the `http` port of the kubernetes service.

```
telepresence intercept svcmid --port 8080:http
```

This will produce output like:

```
Using Deployment svcmid
intercepted
   Intercept name         : svcmid
   State                  : ACTIVE
   Workload kind          : Deployment
   Destination            : 127.0.0.1:8080
   Service Port Identifier: http
   Volume Mount Error     : sshfs is not installed on your local machine
   Intercepting           : all TCP requests
```

See https://www.getambassador.io/docs/telepresence/latest/howtos/intercepts.

### Run the Middle Service Locally

Now that there is an intercept configured, the middle service needs to be started locally.

Edit `svcmid/main.go` to simulate a change in the middle service. Update the `call` endpoint to use a new handler:

```
r.GET("/call", handler2)
```

Add a new `handler2` func to `main.go` with the following contents:

```
func handler2(c *gin.Context) {
	url := "https://scottsbaldwin.github.io/weatherapi/weather/austin"
	res, _ := http.Get(url)
	body, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	var weatherResponse WeatherResponse
	var svcResponse ServiceResponse
	if err := json.Unmarshal(body, &weatherResponse); err != nil {
		svcResponse.Error = err.Error()
		c.JSON(http.StatusBadRequest, svcResponse)
	} else {
		svcResponse.Message = fmt.Sprintf("Current temperature is %.2f°F, but it feels like %.2f°F.", weatherResponse.Current.TempF, weatherResponse.Current.FeelsLike)
		c.JSON(http.StatusBadRequest, svcResponse)
	}
}
```

Start up a local copy of the service:

```
cd svcmid
make run_local
```

The Gin process will start up locally and listen on port 8080:

```
[GIN-debug] Listening and serving HTTP on :8080
```

Run the following command to see that `svctop` handles the request in kuberentes but then call `svcmid`'s `call` endpoint using the intercepted service running outside the cluster.

```
curl http://svctop.default:8080/call
```

The response should now look like (which came from `svcmid` running locally):

```
{"message":"Current temperature is 71.22°F, but it feels like 71.85°F."}
```

### "Leave" the Intercept

To shutdown the intercept and use `svcmid` from the cluster instead of the local copy, run the following command:

```
telepresence leave svcmid
```

Confirm that `svctop` is now consuming the in-cluster `svcmid` service:

```
curl http://svctop.default:8080/call
{"message":"Hi, I am svcbot!"}
```

## References

- https://www.telepresence.io/docs/latest/quick-start/