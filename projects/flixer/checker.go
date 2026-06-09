// ///////////////////////////////////////
// checker.go - Upcoming title checker
// Mike Schilli, 2026 (m@perlmeister.com)
// ///////////////////////////////////////
package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/playwright-community/playwright-go"
)

type upcomingCheck struct {
	Title string
	URL   string
}

func CheckUpcoming(picks []Pick) error {
	checks := []upcomingCheck{}

	for _, idx := range UpcomingPickIndexes(picks) {
		checks = append(checks, upcomingCheck{
			Title: picks[idx].Title,
			URL:   picks[idx].URL,
		})
	}

	if len(checks) == 0 {
		fmt.Println("No upcoming titles found.")
		return nil
	}

	pw, err := playwright.Run()
	if err != nil {
		return err
	}
	defer pw.Stop()

	browser, err := pw.Firefox.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
		SlowMo:   playwright.Float(100),
	})
	if err != nil {
		return err
	}
	defer browser.Close()

	page, err := browser.NewPage()
	if err != nil {
		return err
	}

	for _, check := range checks {
		status, err := checkNetflixTitle(page, check)
		if err != nil {
			fmt.Printf("ERROR\t%s\t%s\t%s\n", check.Title, check.URL, err)
			continue
		}
		fmt.Printf("%s\t%s\t%s\n", status, check.Title, check.URL)
	}

	return nil
}

func checkNetflixTitle(page playwright.Page, check upcomingCheck) (string, error) {
	_, err := page.Goto(check.URL, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
		Timeout:   playwright.Float(60000),
	})
	if err != nil {
		return "", err
	}

	_ = page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State:   playwright.LoadStateNetworkidle,
		Timeout: playwright.Float(15000),
	})
	page.WaitForTimeout(1500)

	pageTitle, _ := page.Title()
	body, _ := page.Locator("body").InnerText(playwright.LocatorInnerTextOptions{
		Timeout: playwright.Float(10000),
	})

	candidates := []string{}
	values, err := page.Locator("a, button, [role=link], h1, h2, h3, [data-uia*='title']").EvaluateAll(`els => els.map(el =>
		el.innerText || el.getAttribute("aria-label") || el.textContent || ""
	)`)
	if err == nil {
		if textValues, ok := values.([]interface{}); ok {
			for _, value := range textValues {
				if text, ok := value.(string); ok {
					candidates = append(candidates, text)
				}
			}
		}
	}

	return classifyNetflixPage(cleanUpcomingTitle(check.Title), pageTitle, candidates, body), nil
}

func classifyNetflixPage(expected, pageTitle string, candidates []string, body string) string {
	lines := titleLines(append([]string{pageTitle}, candidates...))
	for _, line := range lines {
		if line == expected ||
			line == "watch "+expected ||
			strings.HasPrefix(line, "watch "+expected+" ") ||
			strings.HasPrefix(line, expected+" netflix") {
			return "FOUND"
		}
	}

	bodyText := normalizeText(body)
	similarNeedles := []string{
		"similar titles",
		"explore titles related",
		"titles related to",
		"more like this",
		"we don t have",
		"results for",
	}

	for _, needle := range similarNeedles {
		if strings.Contains(bodyText, needle) {
			return "SIMILAR"
		}
	}

	return "UNKNOWN"
}

func cleanUpcomingTitle(value string) string {
	return normalizeText(strings.ReplaceAll(value, "(upcoming)", " "))
}

func titleLines(values []string) []string {
	lines := []string{}

	for _, value := range values {
		for _, line := range strings.Split(value, "\n") {
			line = normalizeText(line)
			if line != "" {
				lines = append(lines, line)
			}
		}
	}

	return lines
}

var nonAlphaNum = regexp.MustCompile(`[^a-z0-9]+`)

func normalizeText(value string) string {
	value = strings.ToLower(value)
	value = nonAlphaNum.ReplaceAllString(value, " ")
	return strings.Join(strings.Fields(value), " ")
}
