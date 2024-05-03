# mathSheets

Introduction:

## SSL to PSQL docker container
1. First use docker command to enter bash shell inside container 
	$ docker exec -it contianerID /bin/sh  
2. Second enter command to go into DB with user:user and password: password 
	$ psql -h 0.0.0.0 -p 5432 -d post_database -U user_post --password


## models

this one is one of the best open AI models:  "text-davinci-002" 

# API Documentation
### Overview
This API allows users to generate math homework problems using the OpenAI API.

# Base URL
The base URL for all endpoints is api/mathsheets.

# Authentication
Authentication is required for certain endpoints. Users must sign up and sign in to obtain authentication tokens.

## Endpoints
1. Sign Up
- URL: /user/signup
- Method: POST
- Description: Creates a new user account.
- Request Body:
``` json
{
    "username": "string",
    "password": "string"
}
```
Response:
- 200 OK if successful, returns authentication token.
- 400 Bad Request if invalid request format.
- 409 Conflict if username already exists.

2. Sign In
- URL: /user/signin
- Method: POST
- Description: Authenticates an existing user.
Request Body:
```json
{
    "username": "string",
    "password": "string"
}
```
Response:
- 200 OK if successful, returns authentication token.
- 400 Bad Request if invalid request format.
- 401 Unauthorized if invalid credentials.

3. Sign Out
- URL: /user/signout
- Method: POST
- Description: Invalidates user's authentication token.
- Authentication: Required
Response:
- 204 No Content if successful.
- 401 Unauthorized if invalid or expired token.

4. User Information
- URL: /user/info
- Method: GET
- Description: Retrieves user information.
- Authentication: Required
- Response:
```json
{
    "username": "string",
    "email": "string",
    "created_at": "timestamp"
}
```

5. Generate Worksheet
- URL: /user/worksheet
- Method: POST
- Description: Generates a math worksheet based on the provided query.
- Request Body:
```json
{
    "query": "string"
}
```
- Authentication: Required
Response:
- 200 OK if successful, returns a PDF file containing the worksheet.
- 400 Bad Request if invalid request format or query.
Error Responses
- 400 Bad Request: Invalid request format or missing parameters.
- 401 Unauthorized: Authentication failure.
- 404 Not Found: Endpoint not found.
- 409 Conflict: Resource conflict (e.g., username already exists).
- 500 Internal Server Error: Server encountered an unexpected condition.
