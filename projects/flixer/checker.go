// ///////////////////////////////////////
// checker.go - Upcoming title checker
// Mike Schilli, 2026 (m@perlmeister.com)
// ///////////////////////////////////////
package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/playwright-community/playwright-go"
)

type upcomingCheck struct {
	Title string
	URL   string
}

type upcomingCheckResult struct {
	Status       string
	JustWatchURL string
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
		result, err := checkNetflixTitle(page, check)
		if err != nil {
			fmt.Printf("ERROR\t%s\t%s\t%s\t%s\n", check.Title, check.URL, result.JustWatchURL, err)
			continue
		}
		fmt.Printf("%s\t%s\t%s\t%s\n", result.Status, check.Title, check.URL, result.JustWatchURL)
	}

	return nil
}

func checkNetflixTitle(page playwright.Page, check upcomingCheck) (upcomingCheckResult, error) {
	title := cleanUpcomingTitle(check.Title)
	searchURL := JustWatchSearchURL(check.Title)
	result := upcomingCheckResult{
		JustWatchURL: searchURL,
	}

	_, err := page.Goto(searchURL, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
		Timeout:   playwright.Float(60000),
	})
	if err != nil {
		return result, err
	}

	waitForPage(page)

	candidates, err := justWatchCandidates(page)
	if err != nil {
		return result, err
	}

	href, ok := findJustWatchCandidate(title, candidates)
	if !ok {
		result.Status = "SIMILAR"
		return result, nil
	}
	result.JustWatchURL = href

	_, err = page.Goto(href, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
		Timeout:   playwright.Float(60000),
	})
	if err != nil {
		return result, err
	}

	waitForPage(page)

	body, _ := page.Locator("body").InnerText(playwright.LocatorInnerTextOptions{
		Timeout: playwright.Float(10000),
	})
	labels, _ := page.Locator("[alt], [aria-label], a, button").EvaluateAll(`els => els.map(el =>
		el.innerText || el.getAttribute("aria-label") || el.getAttribute("alt") || el.textContent || ""
	)`)

	text := body + "\n" + strings.Join(stringValues(labels), "\n")
	if justWatchHasNetflixOffer(text) {
		result.Status = "FOUND"
		return result, nil
	}

	result.Status = "UNKNOWN"
	return result, nil
}

func waitForPage(page playwright.Page) {
	_ = page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State:   playwright.LoadStateNetworkidle,
		Timeout: playwright.Float(15000),
	})
	page.WaitForTimeout(1500)
}

type justWatchCandidate struct {
	Title string
	Href  string
}

func justWatchCandidates(page playwright.Page) ([]justWatchCandidate, error) {
	values, err := page.Locator("a[href*='/movie/'], a[href*='/tv-show/']").EvaluateAll(`els => els.map(el => ({
		title: el.innerText || el.getAttribute("aria-label") || el.textContent || "",
		href: el.href || ""
	}))`)
	if err != nil {
		return nil, err
	}

	candidates := []justWatchCandidate{}
	for _, value := range objectValues(values) {
		title, _ := value["title"].(string)
		href, _ := value["href"].(string)
		if title == "" || href == "" {
			continue
		}
		candidates = append(candidates, justWatchCandidate{Title: title, Href: href})
	}

	return candidates, nil
}

func findJustWatchCandidate(expected string, candidates []justWatchCandidate) (string, bool) {
	for _, candidate := range candidates {
		for _, line := range titleLines([]string{candidate.Title}) {
			if line == expected {
				return candidate.Href, true
			}
		}
	}

	for _, candidate := range candidates {
		if justWatchSlug(candidate.Href) == expected {
			return candidate.Href, true
		}
	}

	return "", false
}

func justWatchSlug(href string) string {
	u, err := url.Parse(href)
	if err != nil {
		return ""
	}

	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) == 0 {
		return ""
	}

	return normalizeText(strings.ReplaceAll(parts[len(parts)-1], "-", " "))
}

func justWatchHasNetflixOffer(value string) bool {
	text := normalizeText(value)
	needles := []string{
		"netflix",
		"netflix basic with ads",
		"netflix standard with ads",
	}

	for _, needle := range needles {
		if strings.Contains(text, needle) {
			return true
		}
	}

	return false
}

func stringValues(values interface{}) []string {
	stringsOut := []string{}

	for _, value := range scalarValues(values) {
		if text, ok := value.(string); ok {
			stringsOut = append(stringsOut, text)
		}
	}

	return stringsOut
}

func objectValues(values interface{}) []map[string]interface{} {
	objects := []map[string]interface{}{}

	for _, value := range scalarValues(values) {
		if object, ok := value.(map[string]interface{}); ok {
			objects = append(objects, object)
		}
	}

	return objects
}

func scalarValues(values interface{}) []interface{} {
	if values == nil {
		return nil
	}

	if result, ok := values.([]interface{}); ok {
		return result
	}

	return nil
}

func cleanUpcomingTitle(value string) string {
	return normalizeText(strings.ReplaceAll(value, "(upcoming)", " "))
}

func JustWatchSearchURL(title string) string {
	return "https://www.justwatch.com/us/search?q=" + url.QueryEscape(cleanUpcomingTitle(title))
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
