
## Your adapter should be able to:

# 1. Parse all possible credential sources (header, form, TLS cert, JWT).
# 2. Capture all grant-type-specific parameters.
# 3. Leave validation (is this combination allowed?) to the domain layer.


HTTP POST /token request combinations:

Grant Type              Auth Method          Example Required Fields
-------------------------------------------------------------------------------------------------------------
authorization_code      Basic Auth           grant_type=authorization_code, code, redirect_uri
authorization_code      Form creds           grant_type=authorization_code, code, redirect_uri, client_id, client_secret
authorization_code      PKCE Public Client   grant_type=authorization_code, code, redirect_uri, code_verifier, client_id
refresh_token           Basic Auth           grant_type=refresh_token, refresh_token
refresh_token           Public client        grant_type=refresh_token, refresh_token, client_id


## Scenarios            Form Values                                 Auth Header           Expected
Public client           client_id=pub123                            -                     ClientID=pub123,ClientSecret=""
Confidential (form)     client_id=conf123&client_secret=shhhh       -                     ClientID=conf123,ClientSecret=shhh
Confidential (basic)    -                                           Basic conf123:shhh    ClientID=conf123,ClientSecret:shhh
Mixed(form wins)        client_id=conf123&client_secret=shhh        Basic other:wrong     ClientID=conf123,ClientSecret:shhh

