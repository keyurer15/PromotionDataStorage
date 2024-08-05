# PromotionDataStorage
The repository stores code files prepared by keyur patel for the storage problem assignment by Verve
The solution is implemented as two services PromotionService and DataStoreService, the implementation is written in Go language (1.22)

**To run the services:**
1) Clone/download the "PromotionService" and "DataStoreService" folders
2) Go inside the PromotionService folder and start CLI ( command line interface)
3) Run "go mod init main"
4) Run "go get github.com/gobwas/ws"
5) Run "go get github.com/gobwas/ws/wsutil"
6) Run "go build"
7) If the build shows no error, run "go run ." to start PromotionService
8) Go inside the DataStoreService folder and start CLI ( command line interface)
9) Run "go mod init main"
10) Run "go get github.com/gobwas/ws"
11) Run "go get github.com/gobwas/ws/wsutil"
12) Run "go build"
13) If the build shows no error, run "go run ." to start DataStoreService

**To test the application**
1) Find a valid promotion ID (like e2649ca5-7e05-4d53-a8ff-919917a4922e) from promotion.csv ( it can be found in the DataStoreService folder)
2) start CLI
3) Run "curl http://localhost:1321/promotions/e2649ca5-7e05-4d53-a8ff-919917a4922e" to send request
   
   Note the protocol is http (not https)
   
4) Observe the response and match it with the corresponding data in the promotions.csv
   
   


   
**Questions**

Q1) The .csv file could be very big (billions of entries) - how would your application perform? How would you optimize it?

A1) When the size of the .csv file is huge (billions of entries), the data file open and read operations become slower, which would introduce latency in the promotion ID search request.
    
    To counter the problem, we should partition the file before the read operation begins, based upon the size of the file, the file can be horizontally partitioned into multiple halves; while creating halves, 
    the new line character can be sought to ensure consistency of records at the boundary of the halves. 
    
    The smaller and more manageable chunks of the data file improve read performance and help achieve future 
    scalability.
    
    The data within the partitions can be analyzed for an indexing strategy, having indexes assigned to partitioned data would improve the search of records across the partitions.    

Q2) How would your application perform in peak periods (millions of requests per minute)? How would you optimize it?

A2) During peak periods (millions of requests per minute) the demand for CPU, network and IO increases multifold, having a single instance of services would not be sufficient to handle a high load of requests.
      
      To serve millions of requests per minute
      
      1)	We should have multiple instances of services to share request traffic
      
      2)	We should introduce a load balancer to distribute traffic efficiently
      
      3)	We can add a cache layer(LRU cache) for quicker access to frequently requested data, caching will also help to reduce the number of requests handled by services
      
      4)	We can analyze the bottleneck junctures and adapt them to non-blocking behaviours (async calls / event-driven, use of web sockets)

Q3) How would you operate this app in production (e.g. deployment, scaling, monitoring)?

A3) To operate this application in production   
    
    The application can be improved to adapt to REST API based architecture, this could help for application operability and maintainability considering the expansion of functionalities.
    
    The application can also have IAM (Identity And Access management) included for authorized access.
    
    The application should be well integrated with continuous integration and deployment ( GitOps way preferably)
   
      a.	Deployment & Scaling
      
            The application can be containerized and hosted on orchestration platforms (cloud based), having such deployment improves the scalability, availability and security aspects of the application.

            Data sharding and replication should be realized, and LRU cache can be deployed to improve performance.
      
      b.	Monitoring
            The deployment can have Splunk and/or Grafana integration for monitoring and alerting purposes, the health of the application can be monitored through Splunk dashboards, and error / abnormal behaviours can be reported through alerts




   
