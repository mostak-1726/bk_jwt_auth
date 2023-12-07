This package build on JWT authentication plugin and echo framework. It exposes two method - GenerateAuthToken and VerifyAuthToken
GenerateAuthToken method verify requested userName and password against configured username and password and return an id_token 
VerifyAuthToken method verify the generated previously and provide jwtToken for further communication


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
and then handler method in router as follows - 
e := echo.New()
e.POST("/bkash/auth", handler.GenerateAuthToken)
e.POST("/bkash/auth/verify", handler.VerifyAuthToken)
Note: Don't forget to add those route to jwt authSkipper 