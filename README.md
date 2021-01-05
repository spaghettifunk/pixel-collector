# Pixel Collector

Simple service that collects a payload and pushes the message into the Bus

## How to develop

Assuming you have docker set it up, run `docker-compose up -d` to launch all the containers. You can then run the api with `go run main.go collect`.

### Tests

You can launch the testing suite with `go clean -testcache && go test -race ./...`. You can add `-v` on the end of the command if you want to verbose all the tests.

You can see the events popping up in Kafka. Exec into the kafka container with `docker exec -it kafka bash` then `cd /usr/bin` and run the following command `./kafka-console-consumer --bootstrap-server kafka:29092 --topic events.raw --from-beginning`.

### Metrics

There are some default metrics exposed with `prometheus` and some custom. To see the them open up your browser to `http://localhost:9090/metrics`. Refresh to have the most updated metrics. The default refresh interval is `5s`. Meaning that every 5s Prometheus pushes the metrics on that endpoint.

## Example of URL

`http://localhost:8082/collect?id=R29X8&uid=1-ovbam3yz-iolwx617&ev=pageload&ed=&v=1&dl=http://edgartownharbor.com/&rl=&ts=1464811823300&de=UTF-8&sr=1680x1050&vp=874x952&cd=24&dt=Edgartown%20Harbormaster&bn=Chrome%2050&md=false&ua=Mozilla/5.0%20(Macintosh;%20Intel%20Mac%20OS%20X%2010_11_5)%20AppleWebKit/537.36%20(KHTML,%20like%20Gecko)%20Chrome/50.0.2661.102%20Safari/537.36&utm_source=&utm_medium=&utm_term=&utm_content=&utm_campaign=`
