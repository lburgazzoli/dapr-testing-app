# Helm chart to install the testing application

To Install the chart, this chart also install Redis (Standalone mode): 

```
helm install testing-app oci://docker.io/salaboy/testing-app --version v0.1.0
```
To interact with the `testing-app-service`:

```
kubectl  port-forward svc/testing-app-service 8080:80
```

Then to read information you can send a GET request to the /read endpoint:
```
http :8080/read
```
To write data to the state store you can send a POST request: 

```
http POST ":8080/write?message=hello"
```

For more info check: [https://www.salaboy.com](https://www.salaboy.com)
