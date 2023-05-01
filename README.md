## docker events

This is a simple go app that collects metrics from container status codes via docker events api.


### how to run?

#### mac

```
 docker run -p 8080:8080 -v /var/run/docker.sock:/var/run/docker.sock:ro docker-events:latest
```

```
go run main.go
http://localhost:8080/metrics
```

### docker exit codes

```
125: docker run itself fails
126: contained command cannot be invoked
127: if contained command cannot be found
128 + n Fatal error signal n:
130 = (128+2) Container terminated by Control-C
137 = (128+9) Container received a SIGKILL
143 = (128+15) Container received a SIGTERM
```