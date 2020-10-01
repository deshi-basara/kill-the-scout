package scraper

import (
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"

	"github.com/deshi-basara/kill-the-scout/database"
	"github.com/tebeka/selenium"
)

type Site struct {
	Name             string
	URL              string
	SearchURL        string
	ItemSelector     string
	ItemLinkSelector string
	TitleSelector    string
	AddressSelector  string
	RentSelector     string
	SizeSelector     string
	RoomsSelector    string
}

type SiteItem struct {
	SiteID string
	ID     string
	URL    string
}

var SitesImmo = Site{
	Name:            "immo",
	ItemSelector:    ".result-list__listing",
	TitleSelector:   ".criteriagroup h1",
	AddressSelector: ".is24-expose-address",
	RentSelector:    ".is24qa-kaltmiete.is24-value",
	SizeSelector:    ".is24qa-flaeche",
	RoomsSelector:   ".is24qa-zi.is24-value",
}

const (
	seleniumPath    = "_vendor/selenium-server.jar"
	geckoDriverPath = "_vendor/geckodriver"
	port            = 4040
)

func ScrapeResultListWithBrowser(url string) []string {
	// Start a Selenium WebDriver server instance (if one is not already running)
	opts := []selenium.ServiceOption{
		selenium.StartFrameBuffer(),           // Start an X frame buffer for the browser to run in.
		selenium.GeckoDriver(geckoDriverPath), // Specify the path to GeckoDriver in order to use Firefox.
		// selenium.Output(os.Stderr),            // Output debug information to STDERR.
	}
	// selenium.SetDebug(true)
	service, err := selenium.NewSeleniumService(seleniumPath, port, opts...)
	if err != nil {
		panic(err) // panic is used only as an example and is not otherwise recommended.
	}
	defer service.Stop()

	// Connect to the WebDriver instance running locally.
	caps := selenium.Capabilities{"browserName": "firefox", "intl.accept_languages": "de,de_DE"}
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		panic(err)
	}
	defer wd.Quit()

	// Navigate to result list
	log.Infoln("Opening: ", url)
	if err := wd.Get(url); err != nil {
		panic(err)
	}

	spew.Dump(wd.Title())

	if err := wd.Wait(Enabled(selenium.ByCSSSelector, SitesImmo.ItemSelector)); err != nil {
		panic(err)
	}

	// Get a reference to a link
	items, err := wd.FindElements(selenium.ByCSSSelector, SitesImmo.ItemSelector)
	if err != nil {
		panic(err)
	}

	results := []string{}
	for _, item := range items {
		link, err := item.FindElement(selenium.ByCSSSelector, "a")
		if err != nil {
			panic(err)
		}
		href, err := link.GetAttribute("href")
		if err != nil {
			panic(err)
		}

		exposeURL := href

		results = append(results, exposeURL)
	}

	return results
}

func ScrapeExposeWithBrowser(exposeID string) database.Expose {
	// build expose-url
	exposeURL := "https://www.immobilienscout24.de/" + exposeID
	var expose database.Expose

	log.Info("Scrape expose: ", exposeURL)

	// Start a Selenium WebDriver server instance (if one is not already running)
	opts := []selenium.ServiceOption{
		selenium.StartFrameBuffer(),           // Start an X frame buffer for the browser to run in.
		selenium.GeckoDriver(geckoDriverPath), // Specify the path to GeckoDriver in order to use Firefox.
		// selenium.Output(os.Stderr),            // Output debug information to STDERR.
	}
	// selenium.SetDebug(true)
	service, err := selenium.NewSeleniumService(seleniumPath, port, opts...)
	if err != nil {
		panic(err) // panic is used only as an example and is not otherwise recommended.
	}
	defer service.Stop()

	// Connect to the WebDriver instance running locally.
	caps := selenium.Capabilities{"browserName": "firefox"}
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		panic(err)
	}
	defer wd.Quit()

	// navigate to expose
	log.Infoln("Opening: ", exposeURL)
	if err := wd.Get(exposeURL); err != nil {
		panic(err)
	}

	if err := wd.Wait(Enabled(selenium.ByCSSSelector, SitesImmo.TitleSelector)); err != nil {
		panic(err)
	}

	// get a reference to title-element
	titleElem, err := wd.FindElement(selenium.ByCSSSelector, SitesImmo.TitleSelector)
	if err != nil {
		panic(err)
	}
	title, err := titleElem.Text()
	if err != nil {
		panic(err)
	}

	log.Infoln("title", title)

	// fet a reference to address-element
	addressElem, err := wd.FindElement(selenium.ByCSSSelector, SitesImmo.AddressSelector)
	if err != nil {
		panic(err)
	}
	address, err := addressElem.Text()
	if err != nil {
		panic(err)
	}

	log.Infoln("address", address)

	// build expose
	expose = database.Expose{
		ID:        exposeID,
		URL:       exposeURL,
		Title:     title,
		Address:   address,
		Rent:      "rent",
		Size:      "size",
		Rooms:     "rooms",
		CreatedAt: time.Now(),
	}

	return expose
}

func Enabled(by, elementName string) func(selenium.WebDriver) (bool, error) {
	return func(wd selenium.WebDriver) (bool, error) {
		el, err := wd.FindElement(by, elementName)
		if err != nil {
			return false, nil
		}
		enabled, err := el.IsEnabled()
		if err != nil {
			return false, nil
		}

		if !enabled {
			return false, nil
		}

		return true, nil
	}
}
