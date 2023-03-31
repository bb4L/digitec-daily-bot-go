// Handling and running the task to query the current daily offer.
package task

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/bb4L/digitec-daily-bot-go/storage"
	"github.com/gocolly/colly"
	"github.com/robfig/cron/v3"
)

var logger = log.New(os.Stdout, "[task] ", log.Ldate|log.Ltime|log.Lmsgprefix)
var dailyOfferChannel chan string

// default values will be overridden with [SetupValues]
var OFFFER_TEMPLATE_STRING = ""
var URL = ""

type dailyOffer struct {
	ItemName         string
	PriceInformation string
	URL              string
}

// SetupValues sets up somem values needed for the task to be performed
func SetupValues(storage *storage.StorageHelper) {
	URL = storage.GetTaskSettings().Url
	OFFFER_TEMPLATE_STRING = storage.GetTaskSettings().CurrentTextTemplate
}

func (offer *dailyOffer) getMessage() string {
	return fmt.Sprintf(OFFFER_TEMPLATE_STRING,
		offer.ItemName, offer.PriceInformation, offer.URL)
}

// StartTaks starts the task to be run cron like at 00:30 for the time "Europe/Zurich"
func StartTask(offerChannel chan string) {
	location, err := time.LoadLocation("Europe/Zurich")
	if err != nil {
		logger.Panicln("error in start task", err)
	}
	c := cron.New(cron.WithSeconds(), cron.WithChain(cron.Recover(cron.DefaultLogger)), cron.WithLocation(location))
	c.AddFunc("0 30 0  * * *", runDigitecCron)
	c.Start()

	dailyOfferChannel = offerChannel
	logger.Println("cron started")
	logger.Println(c.Entries())
}

func runDigitecCron() {
	logger.Println("run digitec cron")
	offer, e := parseDailyOffer()
	if e != nil {
		logger.Println("error parsing daily offer", e)
		return
	}
	dailyOfferChannel <- offer.getMessage()
}

// GetMessageText returns the current daily offer.
func GetMessageText(storage *storage.StorageHelper) string {
	offer, e := parseDailyOffer()
	if e != nil {
		logger.Println("error parsing daily offer", e)
		return ""
	}
	return offer.getMessage()
}

func parseDailyOffer() (dailyOffer, error) {
	var offer = dailyOffer{"", "", ""}

	c := colly.NewCollector()

	c.OnHTML("article.sc-v6swez-1.cpZwnD", func(e *colly.HTMLElement) {
		e.ForEach(
			"a", func(i int, h *colly.HTMLElement) {
				href := h.Attr("href")
				if strings.HasPrefix(href, "/en/product/") {
					offer.URL = "https://digitec.ch" + href
				}
			},
		)
	})

	c.OnHTML("div.sc-v6swez-4.gSBxIL", func(e *colly.HTMLElement) {
		e.ForEach(
			"p", func(i int, h *colly.HTMLElement) {
				if i == 0 {
					offer.ItemName = h.Text
				}
			},
		)

		e.ForEach(
			"div", func(i int, h *colly.HTMLElement) {
				if strings.Contains(h.Text, "was") {
					priceInformation := strings.Replace(h.Text, "was", " was", -1)
					offer.PriceInformation = priceInformation
				}
			},
		)

	})

	c.OnRequest(func(r *colly.Request) {
		logger.Println("Visiting", r.URL)
	})

	c.Visit(URL)

	logger.Println("found offer ", offer)
	return offer, nil
}
