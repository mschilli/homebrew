// ///////////////////////////////////////
// outfile.go - output file determination
// Mike Schilli, 2020 (m@perlmeister.com)
// ///////////////////////////////////////
package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

func outfile(infiles []string) string {
	if len(infiles) == 0 {
		panic("Cannot have zero infiles")
	}

	ext := filepath.Ext(infiles[0])
	base := longestSubstr(infiles)
	base = strings.TrimSuffix(base, ext)
	base = strings.TrimSuffix(base, "-")

	return fmt.Sprintf(
		"%s-out%s", base, ext)
}

func longestSubstr(all []string) string {
	testIdx := 0
	keepGoing := true

	for keepGoing {
		var c byte

		for _, instring := range all {
			if testIdx >= len(instring) {
				keepGoing = false
				break
			}

			if c == 0 { // uninitialized?
				c = instring[testIdx]
				continue
			}

			if instring[testIdx] != c {
				keepGoing = false
				break
			}

		}
		testIdx++
	}

	if testIdx <= 1 {
		return ""
	}
	return all[0][0 : testIdx-1]
}
