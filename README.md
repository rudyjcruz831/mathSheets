# mathSheets

Introduction:

## SSL to PSQL docker container
1. First use docker command to enter bash shell inside container 
	$ docker exec -it contianerID /bin/sh  
2. Second enter command to go into DB with user:user and password: password 
	$ psql -h 0.0.0.0 -p 5432 -d post_database -U user_post --password