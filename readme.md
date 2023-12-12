# Echo JWT Auth Integrator
This package builds on JWT authentication plugin and echo framework. It exposes two methods - GenerateAuthToken and VerifyAuthToken.
GenerateAuthToken method verifies the requested username and password against the configured username and password and returns an id_token. 
VerifyAuthToken method verifies the token generated previously and provides jwtToken for further communication

## Installation
To Install run the following command -  
```bash
go get github.com/mostak-1726/bk_jwt_auth
```

## How to use:
Instantiate capAuthIntegrator as follows -
```go
import cap_auth "github.com/mostak-1726/bk_jwt_auth"
import "github.com/go-redis/redis"
c := cap_auth.Config{
        UserName:             "abc",
        Password:             "12345",
        ExpiryInSec:          3600,
        TestCustomerAppToken: "14580760-b5d9-42d7-aa3a-51d20caeff6a",
        JwtTokenSecret:      "testSecret",
        RedisClient:          *redis.Client,
}
handler := auth.NewCapAuthIntegrator(c)
e := echo.New()
e.POST("/auth", handler.GenerateAuthToken)
e.POST("/auth/verify", handler.VerifyAuthToken)
```
/auth endpoint needs to request with the following params - 
```json

{
    "username": "abc",
    "password": "12345",
    "mobile_number": "01799021432"
}
```

/auth/verify endpoint needs to request with the following params - 
```json
{
    "token": "6c843576-5ac0-4ba9-8f44-6265d3e18039",
    "mobile_number": "01799021432"
}
```
Note: Don't forget to add those routes to jwt authSkipper
