// ///////////////////////////////////////
// args.go - pdftk argument collector
// Mike Schilli, 2020 (m@perlmeister.com)
// ///////////////////////////////////////
package main

import "fmt"

func pdftkArgs(files []string) []string {
	args := []string{"pdftk"}
	catArgs := []string{}
	letterChr := int('A')

	for idx, file := range files {
		letter := string(letterChr + idx)
		args = append(args,
			fmt.Sprintf("%s=%s", letter, file))
		catArgs = append(catArgs,
			fmt.Sprintf("%s1-end", letter))
	}

	args = append(args, "cat")
	args = append(args, catArgs...)
	args = append(args,
		"output", outfile(files))
	return args
}
