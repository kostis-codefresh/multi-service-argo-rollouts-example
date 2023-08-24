# Simple notification tester

This is an example application to be used with the [Codefresh GitOps Certification](https://learning.codefresh.io/)

It is simple application that consists of one container.


It accepts webhooks of kind at `/notify` and then shows them 
in a simple list view

![List-view](list-view.png)

## How to run locally

`go run .`

then visit http://localhost:8080 in your browser

## How to build and run a container

Run

 *  `docker build . -t my-app` to create a container image 
 *  `docker run -p 8080:8080 my-app` to run it

 then visit http://localhost:8080 in your browser

You can find prebuilt images at [https://hub.docker.com/r/kostiscodefresh/summer-of-k8s-app/tags](https://hub.docker.com/r/kostiscodefresh/summer-of-k8s-app/tags)




