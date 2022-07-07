This module implements two RESTful API endpoints to consume payloads passed in YAML format.

# /metadata/postMetadata
- the user passes the payload (as YAML) one at a time.


# /metadata/getMetadata: 
- If the key(s) exist, and value(s) match, then return the matched payloads in memory
- If no match, then return all payloads in memory