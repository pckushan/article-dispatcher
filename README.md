# Article Dispatcher  v1.0.0
This is a simple Go web service to dispatch articles using 3 endpoints.

## Getting Started
Prerequisites
Before you can build and run this web service, you need to have the following software installed on your system:

* Go ( version> 1.17 )
  * [install or update](https://go.dev/doc/install)
* Docker

## Building Locally
To build the web service locally, follow these steps:

1. Clone the repository: 
    ```shell
    git clone https://github.com/pckushan/article-dispatcher.git
    ```
2. Change into the project directory: cd articel-dispatcher 
    ```shell
    cd articel-dispatcher 
    ```
3. Build the binary:
    ```shell
    go build -o articel-dispatcher
    ```
   OR in a mac,
    ```shell
    env GOOS=darwin GOARCH=amd64 go build -o articel-dispatcher
    ```
   OR in a linux,
    ```shell
    env GOOS=linux GOARCH=amd64 go build -o articel-dispatcher
    ```

## Running Locally
To run the web service locally, follow these steps:

Make sure you're in the project directory: 
```shell
    cd articel-dispatcher
```
  
Run the binary: 
```shell
./articel-dispatcher
```
The web service will now be running on http://localhost:8888

## Building and Running with Docker
To build and run the web service with Docker, follow these steps:

Clone the repository:
```shell
git clone https://github.com/pckushan/article-dispatcher.git
```
Change into the project directory:
```shell
cd article-dispatcher
```
Build the Docker image:
```shell
docker build -t article-dispatcher .
```
Run the Docker container: 
```shell
docker run -p 8888:8888 article-dispatcher
```
The web service will now be running on http://localhost:8888

## Endpoints 
API documentation can be found [here](docs/swagger.yaml)
This web service has the following endpoints:

POST /articles
Reverses the string passed in the request body.

Example:
```shell
curl --location --request POST 'localhost:8888/articles' \
--data-raw '{
  "id": "1",
  "title": "latest science shows that potato chips are better for you than sugar",
  "date" : "2016-09-23",
  "body" : "some text, potentially containing simple markup about how potato chips are great",
  "tags" : [ "nature", "fitness"]
}'
```
```json
{
  "data": {
    "id" :"1"
  }
}
```

GET /articles/{id}
Returns the corresponding article related with the id.

Example:

```shell
curl --location --request GET 'localhost:8888/articles/1'
```
```json
{
  "id": "1",
  "title": "latest science shows that potato chips are better for you than sugar",
  "date" : "2016-09-23",
  "body" : "some text, potentially containing simple markup about how potato chips are great",
  "tags" : [ "nature", "fitness"]
}
```

GET /tags/{tagName}/{date}
Filters the articles data with the tag related to the date.

Example:

```shell
curl --location --request GET 'localhost:8888/tags/nature/20160923'
```

```json
{
    "tag": "nature",
    "count": 2,
    "articles": [
        "1"
    ],
    "related_tags": [
        "fitness"
    ]
}
```

## Assumptions

- In the last endpoint's implementation, as per the example it showed count 
as 17 and I think it should be 3. Since requirement was to get the count 
of the distinct tags related to the date and tag requested.    

- Requirement was to keep the data in memory and the data is cached to remain until the project 
is up and running. It will NOT persist any data added once the service is restarted.

## Limitations & Improvements

#### request rate limiter

- for large number of requests, need to implement an asynchronous queue
- could use an api-gateway

#### authentication and authorization

- inbuilt middleware can be implemented (such as JWT)
- could use an api-gateway

#### scaling

- can use a kubernetes configuration to scale with multiple instances.

  NOTE: if multiple instances used need to use a central caching mechanism
  or background sync the records to all the instances.