package main

import (
	"os"

	"github.com/asdine/storm"
	"github.com/bamzi/jobrunner"
	"github.com/deshi-basara/kill-the-scout/database"
	"github.com/deshi-basara/kill-the-scout/notifier"
	"github.com/deshi-basara/kill-the-scout/scraper"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	interval := os.Getenv("IMMO_INTERVAL_SECONDS")

	log.Infoln("Kill The Scout started. Checking every", interval)

	jobrunner.Start()
	jobrunner.Schedule("@every "+interval, KillTheScout{})
	jobrunner.Now(KillTheScout{})

	select {}
}

// KillTheScout job functions
type KillTheScout struct {
}

// Run will get triggered automatically.
func (e KillTheScout) Run() {
	url := os.Getenv("IMMO_SEARCH_URL")
	mailTo := os.Getenv("MAIL_TO")

	// open embedded database
	db := database.Setup()

	//  get all latest (at max 20) Exposes from the result-list and check if they were already scraped
	results := scraper.ScrapeResultList(url)

	for _, exposeID := range results {
		_, err := database.GetExpose(db, exposeID)
		if err == storm.ErrNotFound {
			log.Infoln("New expose found:", exposeID)

			// scarpe and save expose
			expose := scraper.ScrapeExpose(exposeID)
			err := database.SaveExpose(db, expose)
			if err != nil {
				log.Fatalf("SaveExpose.Fatal: %s %s \n", err, err == storm.ErrNotFound)
			}

			// notify
			notifier.SendMail(mailTo, expose)
		} else if err != nil {
			log.Fatalf("GetExpose.Fatal: %s \n", err)
		}
	}

	defer db.Close()
}
