This package builds on JWT authentication plugin and echo framework. It exposes two methods - GenerateAuthToken and VerifyAuthToken
GenerateAuthToken method verifies the requested username and password against the configured username and password and returns an id_token 
VerifyAuthToken method verifies the generated previously and provides jwtToken for further communication


To Install run the following command -  
go get github.com/mostak-1726/bk_jwt_auth

How to use:
Instantiate capAuthIntegrator as follows -
conf := _type.RedisConfig{
        Host: "127.0.0.1",
        Port: "6379",
        Pass: "secret_redis",
        Db:   1,
        Ttl:  3600,
   }

c := _type.Config{
        UserName:             "mostak",
        Password:             "12345",
        ExpiryInSec:          3600,
        TestCustomerAppToken: "14580760-b5d9-42d7-aa3a-51d20caeff6a",
        JwtTokenSecrete:      "testSecret",
        RedisConfig:          conf,
}
handler := auth.NewCapAuthIntegrator(c)
and then the handler method in the router as follows - 
e := echo.New()
e.POST("/bkash/auth", handler.GenerateAuthToken)
e.POST("/bkash/auth/verify", handler.VerifyAuthToken)
/bksh/auth endpoint needs to request with following params - 

{
    "username": "mostak",
    "password": "12345",
    "mobile_number": "01799021432"
}
/bkash/auth/verify endpoint needs to request with following params - 
{
    "token": "6c843576-5ac0-4ba9-8f44-6265d3e18039",
    "mobile_number": "01711082738"
}
Note: Don't forget to add those route to jwt authSkipper 
