This module implements two RESTful API endpoints to consume payloads passed in YAML format.

# APIs supported by this application
## /metadata/postMetadata
- the user passes the payload (as YAML) one at a time.


## /metadata/getMetadata: 
- If the key(s) exist, and value(s) match, then return the matched payloads in memory
- If no match, then return all payloads in memory

  
# TESTING
## Using 'curl' from command line:
- TEST postMetadata:   
--  Input valid data  
curl http://localhost:8080/metadata -X POST -H 'Content-Type:application/yaml' --data-binary @valid-1.yml  
curl http://localhost:8080/metadata -X POST -H 'Content-Type:application/yaml' --data-binary @valid-2.yml  

-- Input invalid data:
curl http://localhost:8080/metadata -X POST -H 'Content-Type:application/yaml' --data-binary @invalid-1.yml   
curl http://localhost:8080/metadata -X POST -H 'Content-Type:application/yaml' --data-binary @invalid-2.yml  

- TEST getMetadata:  
-- Retrieve all the data  
curl http://localhost:8080/metadata   

-- Retrieve data based on search parameters:   
curl http://localhost:8080/metadata?title=Valid+App+1\&version=0.0.1   
curl http://localhost:8080/metadata?title=App   


