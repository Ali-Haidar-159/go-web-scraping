package main

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
)

type Card struct {
	Title       string
	Price       string
	Description string
	Image       string
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	ctx, cancel2 := chromedp.NewContext(ctx)
	defer cancel2()

	var cards []Card

	err := chromedp.Run(ctx,
		chromedp.Navigate(`https://www.airbnb.com/s/Kuala-Lumpur/homes?place_id=ChIJ5-rvAcdJzDERfSgcL1uO2fQ&refinement_paths%5B%5D=%2Fhomes&flexible_trip_lengths%5B%5D=weekend_trip&date_picker_type=FLEXIBLE_DATES&search_type=HOMEPAGE_CAROUSEL_CLICK`),
		chromedp.Sleep(7*time.Second),

		// Scroll down to trigger lazy loading of images
		chromedp.Evaluate(`window.scrollTo(0, document.body.scrollHeight)`, nil),
		chromedp.Sleep(3*time.Second),
		chromedp.Evaluate(`window.scrollTo(0, 0)`, nil),
		chromedp.Sleep(1*time.Second),

		chromedp.Evaluate(`
			Array.from(document.querySelectorAll('.g1qv1ctd')).map(card => {

				// Title
				const title = card.querySelector('[data-testid="listing-card-title"]')?.innerText.trim() || "";

				// Description: only grab leaf spans (no children that are also spans)
				// This avoids duplicating parent+child span text
				const subtitleContainer = card.querySelector('[data-testid="listing-card-subtitle"]');
				let desc = "";
				if (subtitleContainer) {
					const leafSpans = Array.from(subtitleContainer.querySelectorAll('span'))
						.filter(s => s.querySelector('span') === null) // leaf only
						.map(s => s.innerText.trim())
						.filter(Boolean);
					// Deduplicate consecutive identical values
					const deduped = leafSpans.filter((v, i, arr) => v !== arr[i - 1]);
					desc = deduped.join(" | ");
				}

				// Price: use the price-availability-row container and grab only direct
				// text from the _first_ matching price span to avoid duplication
				const priceRow = card.querySelector('[data-testid="price-availability-row"]');
				let price = "";
				if (priceRow) {
					// Get all leaf spans, deduplicate
					const priceSpans = Array.from(priceRow.querySelectorAll('span'))
						.filter(s => s.querySelector('span') === null)
						.map(s => s.innerText.trim())
						.filter(Boolean);
					const deduped = [...new Set(priceSpans)];
					price = deduped.join(" ");
				}

				// Image: check multiple attributes for lazy-loaded images
				const imgEl = card.querySelector('img');
				let img = "No Image";
				if (imgEl) {
					img = imgEl.src || imgEl.getAttribute('data-src') || imgEl.getAttribute('data-original') || "No Image";
					// Prefer a loaded srcset entry if src is a placeholder/blank
					if ((!img || img.startsWith('data:') || img === 'No Image') && imgEl.srcset) {
						img = imgEl.srcset.split(',')[0].trim().split(' ')[0];
					}
				}

				return { Title: title, Price: price, Description: desc, Image: img };
			})
		`, &cards),
	)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, data := range cards {
		fmt.Printf("{Title:%s\nPrice:%s\nDescription:%s\nImage:%s}\n", data.Title, data.Price, data.Description, data.Image)
		fmt.Println("--------------------------")
	}
}