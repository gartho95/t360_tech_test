**HOW TO RUN**

Rename ```.env.example``` to ```.env```

IF you have a google project ID with a pubsub set up:

fill in the ```GOOGLE_CLOUD_PROJECT_ID``` in the .env (Given the time constraints I have only briefly tested if this works on a test Google Cloud account)

elsewise keep it empty and the code will return a mocked example


in the project directory run:

```go mod tidy```
```go run main.go```



With more time I would've liked to add unit testing, and proper error handling everywhere as I am aware the way it currently handles it in the concurrency is not good.

Below are the notes I made prior to writing any code, to hopefully give a small idea of my thought process:

## Essentials  
- Get list of test VRMs (vehicle registrations) and store locally in project, seems that the URL containing VRMs calls data server-side, quick web crawler.  
- Create `.env` file to set `BASE_URL` and Google Cloud credentials.  
- Create structs for POST request and responses.  
- Function to build a map of VRMs and their associated company.  
- Function to make requests that handles timeouts.  
- Concurrent handling of requests:  
  - One per data source.  
  - One per VRM for each data source (one concurrent loop for company, another for vehicles within company).  
- Mock function to handle publishing `positive_searches`, once each concurrent process has succeeded.  

## Nice to Have  
- Handle non-200 responses from requests.  
- Full Google Cloud integration.  
- Clean concurrent error handling and logging.  
- Dockerization.  
- GitHub Actions to pull `.env` from staging/production config.  
- Unit tests.  
- Mock batch publishing of `positive_searches`.  
