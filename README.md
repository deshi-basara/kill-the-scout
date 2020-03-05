# kill-the-scout

Tired if pressing F5 all day on the immobilienscout24 search-page? Use this simple script to 
to check immobilienscout24 and get notifications for new exposes/flats via email.

**Note:** KillTheScout will only scrape the last 20 entries from the search-page.

## Configuration

Create an `.env`-file in the root-folder and adjust the settings accordingly:

```
SMTP_SERVER=smtp.gmail.com:587
SMTP_HOST=smtp.gmail.com
SMTP_USER=login@gmail.com
SMTP_PASS=password

IMMO_SEARCH_URL=https://www.immobilienscout24.de/Suche/de/baden-wuerttemberg/stuttgart/wohnung-mieten?price=-4048.0&geocodes=1276001039002&sorting=2&enteredFrom=result_list
IMMO_INTERVAL_SECONDS=180s
MAIL_TO=mail-me-exposes@mail.com
```

## Install & Run

Install `golang`, clone this repo and execute the following commands:

```
> go install
> go run main.go
```
