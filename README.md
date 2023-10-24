# Web Crawling using Golang

This is a backend heavy project which is used to crawl content from different URLs.
I have implemented that when a webpage is accessed for the first time no content is posted but it gets available in the cache. Re-accessing the same webpage results in the content.

In the sendx_assignment folder:
1. The project.go contains the Golang code
2. The project.html file inside the static folder contains the frontend code for a simple web-crawler webpage

## How to run
1. To run the code start the Go server and then enter localhost:8080/static/project.html for the website.
2. To see the timestamps of first access type localhost:8080/view
3. To view the access logs of everytime a webpage is accessed type localhost:8080/accesslogs
