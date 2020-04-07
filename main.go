package main

import (
	"fmt"
	"github.com/deanishe/awgo"
	"github.com/hzenginx/tureng/tureng"
	"golang.org/x/text/unicode/norm"
	"regexp"
	"strings"
)

var wf *aw.Workflow

func init() {
	wf = aw.New()
}

func splitInput(query string) (string, []string) {
	splited := strings.Split(query, ":")
	return splited[0], splited[1:]
}

func getIcon(category string) *aw.Icon {
	var icon string
	matched, _ := regexp.MatchString(`en->tr`, category)
	print(category)
	if matched {
		icon = "en_tr.png"
	} else {
		icon = "tr_en.png"
	}
	return &aw.Icon{
		Value: icon,
	}
}

func handleSearchResponse(response *tureng.SearchResponse, wf *aw.Workflow) {
	if response.IsSuccessful {
		if response.Result.IsFound == 1 {
			for _, result := range response.Result.Results {
				icon := getIcon(result.Category)
				wf.NewItem(result.Term).Subtitle(result.Category).Icon(icon).Arg(result.Term).Valid(true)
			}
		} else {
			for _, suggestion := range response.Result.Suggestions {
				wf.NewWarningItem(suggestion, "Did you mean this?").Autocomplete(fmt.Sprintf("translate:%s", suggestion))
			}
		}
	} else {
		wf.Fatal(response.Exception)
	}
	wf.SendFeedback()
}

func handleAutocompleteResponse(response *tureng.AutoCompleteResponse, wf *aw.Workflow) {
	for _, word := range response.Words {
		wf.NewItem(word).Autocomplete(fmt.Sprintf("translate:%s", word)).Icon(&aw.Icon{Value: "tureng.png"})
	}
	wf.SendFeedback()
}

func run() {
	var query = wf.Args()[0]
	query = norm.NFC.String(query)
	command, args := splitInput(query)

	if command == "translate" {
		if len(args) > 0 {
			word := args[0]
			response, err := tureng.Search(word)
			if err != nil {
				wf.FatalError(err)
			} else {
				handleSearchResponse(response, wf)
			}
		}
	} else {
		response, err := tureng.AutoComplete(command)
		if err != nil {
			wf.FatalError(err)
		} else {
			handleAutocompleteResponse(response, wf)
		}
	}
}

func translate() error {
	wf.NewItem("item")
	wf.SendFeedback()
	return nil
}

func main() {
	wf.Run(run)
}
