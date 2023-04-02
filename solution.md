## Introduction:
In this project, I implemented a simple API that handles article data in JSON format
and allows users to retrieve articles by ID or by tags and dates. The API is 
implemented in Go, a language that excels in **performance** and **concurrency**. I used 
the clean architecture approach to structure the project, dividing  it into several layers. 
Here, I will describe my approach to implementing the API, including my choice of language, 
the architecture, request tracing, error handling, best practices, service monitoring,  and the testing strategy.

## Choice of Language:
- well-suited for **high-performance** and **concurrent** applications.
- has a straightforward syntax that is easy to read and learn. 
- easy to write the code quickly and efficiently. 
- has a robust standard libraries that provides all the tools necessary to implement 
the modules. 

## Architecture:
- Used the clean architecture approach to structure the project.
- This architecture separates the project into, several main layers: 
**domain**, **service**, **http**, and **adaptors**. 
  * The `domain` layer contains the business logic and entities, 
  * the `service` layer implements the use cases and service level processing, 
  * the `http` layer handles the data transport, and 
  * the `adaptors` layer exposes the packages or plugins.

This is not a defined code structure for the clean architecture, but packaged to 
achieve clean architecture. Main idea is to isolate the domain related logics 
and models from the dependency services and plugins and expose them via 
abstractions. This way it can reuse the core domain with other implementations.   

## Request Tracing
A middleware is used to add a `trace id` to all the request coming into the server and 
logs with the request url. This is helpful in finding any errors occur in the system, 
to back track and identify the error.    

```
2023/04/02 22:12:06 [DEBUG]:  request received with trace-id:[2fd3c413-8e83-4a31-ac64-2cf825ce4c34] url [/articles]
```

## Error handling
A http handler was used to handle all the system errors, and it maps to a 
common defined error structure to keep a consistent error message. 
If need to add a new error it needs to map and add to the handler and all the 
logging will happen there at the `error handler` with the trace of the error.   

```json
{
    "code": 40013,
    "description": "invalid article date format",
    "trace": "82412b10-674e-452b-a1f4-3f93bbb349c2"
}
``` 

## Best Practices
* [linter file](.golangci.yml) is used to align with best practices in coding and can validate time 
to time by running, `golangci-lint run`  
NOTE: if command not found, need to [install](https://golangci-lint.run/) `golangci-lint` 

* API contract is created using [swagger](docs/swagger.yaml) 
* [Makefile](Makefile) is used in executing commands 

## Service Monitoring
- service level [matrics](http://localhost:7001/metrics) are integrated (http:localhost:7001/metrics) into the system, and it gives the 
**latencies** and **count** for each endpoint. 
```
nine_article_dispatcher_request_latency_micro_sum{endpoint="add_article",error="false"} 206
nine_article_dispatcher_request_latency_micro_count{endpoint="add_article",error="false"} 1
```
- grafana dashboard need to be integrated to monitor and observe the service performance. 


## Testing Strategy:
-  used the standard Go testing package to write unit tests and integration tests 
for the HTTP API endpoints. The unit tests verify the behavior of the individual 
components, while the integration tests ensure that the components work together 
correctly.
- unit tests were written only for the [`cache` package](internal/adaptors/cache/cache_test.go)
    ```shell
    go test -v -count=1 ./...
    ```
  - [end to end](/e2e-test/e2e_test.go) test was written to sequentially execute 3 endpoints and 
  validate the responses.
  ```shell
      go test -tags=e2e ./e2e-test -v -count=1
    ```
## Conclusion:
In conclusion, I implemented a simple API that handles article data in JSON 
format using Go. I used a modified version of the clean architecture to structure 
the project, dividing it into four layers. The use of this architecture ensured 
that the code was modular, reusable, and maintainable. I also employed a testing 
strategy that included unit and integration tests to ensure that the application
worked correctly. Overall, my solution provides a simple and efficient 
implementation of the required API endpoints while following best practices and 
oooomaintaining code quality.