package codywebapi

import (
	"bufio"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/cody6750/codywebapi/codyWebAPI/tools"
	"github.com/cody6750/codywebapi/codyWebAPI/website/amazon"
)

const (
	amazonFlag         string = "amazon"
	bestbuyFlag        string = "bestbuy"
	websiteFlag        string = "website"
	componentName      string = "codyWebAPI"
	printAction        string = "print"
	searchAction       string = "search"
	minimumInputLength int    = 3
)

var (
	applicationCommand  string
	applicationItem     string
	applicationWebsite  string
	applicationCommands map[string]string = map[string]string{
		"search": "search",
	}
	supportedWebsites map[string]string = map[string]string{
		"amazon": "amazon",
	}
	supportedFlags map[string]struct{} = map[string]struct{}{
		"website": {},
		"item":    {},
	}
	errInput           error = errors.New("Invalid Input")
	errParseFlag       error = errors.New("Error parsing flags")
	errFlagNotSet      error = errors.New("Error, flag not set")
	errUnsupportedFlag error = errors.New("Error, unsupported flag")
	errWebsiteFlag     error = errors.New("Error, unsupported website")
	listOfFlags        []flag
	listInt            []int
)

type flag struct {
	flag      string
	flagValue string
}

//Run ...
func Run() {
	initializeComponent()
}

func initializeComponent() {
	log.Printf("%v Initializing component: %v", tools.FuncName(), componentName)
	log.Printf("%v Finished initializing component %v", tools.FuncName(), componentName)

	scanner := bufio.NewScanner(os.Stdin)
	log.Printf("%v Component %v is up and running. Now waiting for input", tools.FuncName(), componentName)
	for scanner.Scan() {
		parseInput(scanner.Text())
	}
}

func parseInput(input string) error {
	splitInput := strings.Split(input, " ")

	if len(splitInput) < minimumInputLength {
		log.Printf("%v Not enough arguments given", tools.FuncName())
		return errInput
	}
	if splitInput[0] != componentName {
		log.Printf("%v %v is the incorrect program. Please use %v instead", tools.FuncName(), splitInput[0], componentName)
		return errInput
	}
	inputCommand, inputCommandExist := applicationCommands[splitInput[1]]
	if inputCommandExist {
		applicationCommand = inputCommand
	} else {
		log.Printf("%v Command does not exist, please use 'help' to list all available commands", tools.FuncName())
		return errInput
	}
	// Passes all flags from the input
	if parseFlagsError := parseFlags(applicationCommand, splitInput[2:]); parseFlagsError != nil {
		return parseFlagsError
	}
	return nil
}

func parseFlags(command string, input []string) error {
	var currentFlag, flagValue string
	var newFlag flag
	// If the first word isn't a flag, exit early.
	if !strings.HasPrefix(input[0], "--") {
		log.Printf("%v Flag not provided in %v", tools.FuncName(), input)
		return errParseFlag
	}

	// Get flags and assign flags
	for index, word := range input {
		if strings.HasPrefix(word, "--") {
			currentFlag = strings.Trim(word, "--")
			if supported, _ := checkIfFlagIsSupported(currentFlag); supported != true {
				log.Printf("%v "+errUnsupportedFlag.Error(), tools.FuncName())
				return errUnsupportedFlag
			}
			newFlag.flag = currentFlag
			if index < len(input)-1 {
				if strings.HasPrefix(input[index+1], "--") {
					log.Printf("%v "+errFlagNotSet.Error(), tools.FuncName())
					return errFlagNotSet
				}
			}
		} else {
			flagValue = flagValue + " " + word
			if len(input)-1 == index {
				newFlag.flagValue = flagValue
				if newFlag.flag == websiteFlag {
					_, err := checkIfWebsiteIsSupported(flagValue)
					if err != nil {
						return errWebsiteFlag
					}
				}
				listOfFlags = append(listOfFlags, newFlag)
				flagValue = ""
			} else if index < len(input)-1 {
				if strings.HasPrefix(input[index+1], "--") {
					newFlag.flagValue = flagValue
					if newFlag.flag == websiteFlag {
						_, err := checkIfWebsiteIsSupported(flagValue)
						if err != nil {
							return errWebsiteFlag
						}
					}
					listOfFlags = append(listOfFlags, newFlag)
					flagValue = ""
				}
			}
		}
	}
	listOfFlags = nil
	return nil
}

func checkIfFlagIsSupported(flag string) (bool, error) {
	flag = strings.TrimSpace(flag)
	if _, flagSupported := supportedFlags[flag]; flagSupported != true {
		log.Printf("%v "+errFlagNotSet.Error(), tools.FuncName())
		return false, errUnsupportedFlag
	}
	return true, nil
}

func checkIfWebsiteIsSupported(website string) (bool, error) {
	website = strings.TrimSpace(website)
	files, err := ioutil.ReadDir("./codyWebAPI/website")
	if err != nil {
		log.Print(err)
	}
	for _, f := range files {
		if website == f.Name() {
			log.Printf("%v %v is a supported website.", tools.FuncName(), website)
			return true, nil
		}
	}
	log.Printf("%v "+errWebsiteFlag.Error()+" : %v", tools.FuncName(), website)
	return false, errWebsiteFlag
}

func getWebsiteFromFlags() {
	log.Printf("%v Attempting to get website from flags", tools.FuncName())
	switch applicationWebsite {
	case amazonFlag:
		w := amazon.Constructor()
		//t := &amazon.Amazon{}
		invokeAmazonActions(w)
		//invokeAmazonActions(t)
	case bestbuyFlag:
		log.Printf("Not implemented yet")
	default:
		log.Fatalf("%v Failed to complete function, website %s is not supported", tools.FuncName(), applicationWebsite)
	}
}

func invokeAmazonActions(website *amazon.Amazon) {
	switch applicationCommand {
	case printAction:
		website.PrintWebsite()
	case searchAction:
		website.SearchWebsite(applicationItem)
	default:
		log.Fatalf("%v Invalid action,type !help for all avaliable actions", tools.FuncName())
	}
}

func shutDown() {
	log.Printf("Exiting Program")
	os.Exit(1)
}

func test() {

}
