# go-jwt-auth

## TDOO
- [ ] nonce
- [ ] kid

## Login Request

```bash
curl -X POST localhost:5555/login -d "username=test&password=password"
```

```json
{
   "access_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiZXhwIjoxNTk1NjYxNDQwLCJuYW1lIjoiU2FtcGxlIE5hbWUiLCJzdWIiOjF9.5RGP9pnlpZ1EGMLYeyOaVGalHcbkP6zPdvR32GljKk4",
   "refresh_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTU3NDY5NDAsInN1YiI6MX0.dowmE8JsDVCX1zBdILigGiA2CdUaeUWoGOFXBKiCKbQ"
}
```

## RefreshToken Request

```bash
curl -X POST localhost:5555/token -H "Content-Type: application/json" -d '{"refresh_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTU3NTAwNTAsInN1YiI6MX0.1DK_UCrXXbFIy2kvcXsN_kKsBszFoFgor7fhXSq3CbM"}'
```
