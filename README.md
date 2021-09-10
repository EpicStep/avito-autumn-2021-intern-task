# Avito intern task (autumn 2021)
## Run (Locally)
Create .env file and add following values:
```dotenv
PORT=8081
DATABASE_URL=postgresql://postgres@postgres:5432/postgres
EXCHANGERATESAPI_TOKEN=<TOKEN>
```

After you can run app in docker:
```bash
docker compose up -d
```
You need to insert ```migrations/000001_initial.up.sql```