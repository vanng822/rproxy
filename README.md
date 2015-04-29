# rproxy
The idea is to have a proxy where host can be registered dynamically. 
The goal is creating something for development purpose where developers can
register their application server and this proxy will send request to correct place.

# Usage example for fakedomain.tld

## Adding backend node to a server

	curl -d 'serverName=fakedomain.tld&targetUrl=http://127.0.0.1:8080' 'http://127.0.0.1:5556/_server/backend'

## Show list of servers

	curl 'http://127.0.0.1:5556/_server'

## Show list of backend node for a server
	
	curl 'http://127.0.0.1:5556/_server/backend?serverName=fakedomain.tld'

## Delete a backend node

	curl -X DELETE 'http://127.0.0.1:5556/_server/backend?serverName=fakedomain.tld&targetUrl=http://127.0.0.1:8080'
	
