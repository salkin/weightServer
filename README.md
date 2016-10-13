# weightServer
REST server for registering weights.

server uses port 3080 as default can be set in compose file

	Register your weight through RESP API:
	Example: 
	POST /api/v1/weight/Brian 
	{ 
	  "date": "dd/mm/yyyy"  
	  "value": 67.0 
	} 
	
	GET /api/v1/weight/Brian`
