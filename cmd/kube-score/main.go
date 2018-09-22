package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/zegl/kube-score/score"
	"io"
	"log"
	"os"
)

func main() {
	filesToRead := os.Args[1:]
	if len(filesToRead) == 0 {
		log.Println("No files given as arguments")
		os.Exit(1)
 	}

	for _, file := range filesToRead {

		var fp io.Reader

		if file == "-" {
			fp = os.Stdin
		} else {
			var err error
			fp, err = os.Open(file)
			if err != nil {
				panic(err)
			}
		}

		scoreCard, err := score.Score(fp)
		if err != nil {
			panic(err)
		}

		sumGrade := 0

		for _, resourceScores := range scoreCard.Scores {

			firstCard := resourceScores[0]

			color.New(color.FgMagenta).Printf("%s %s in %s\n", firstCard.ResourceRef.Kind, firstCard.ResourceRef.Name, firstCard.ResourceRef.Namespace)

			for _, card := range resourceScores {
				col := color.FgGreen
				status := "OK"

				if card.Grade == 0 {
					col = color.FgRed
					status = "CRITICAL"
				} else if card.Grade < 10 {
					col = color.FgYellow
					status = "WARNING"
				}

				color.New(col).Printf("    [%s] %s\n", status, card.Name)

				for _, comment := range card.Comments {
					fmt.Printf("        * %s\n", comment)
				}

				sumGrade += card.Grade
			}
		}
	}
}
