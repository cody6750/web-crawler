package codywebapi

import (
	"bufio"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/cody6750/codywebapi/codyWebAPI/tools"
	"github.com/cody6750/codywebapi/codyWebAPI/website"
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
	applicationCommands map[string]string = map[string]string{
		"search": "search",
	}
	supportedFlags map[string]struct{} = map[string]struct{}{
		"website": {},
		"item":    {},
	}
	errInput             error = errors.New("invalid Input")
	errParseFlag         error = errors.New("error parsing flags")
	errFlagNotSet        error = errors.New("error, flag not set")
	errUnsupportedAction error = errors.New("error, unsupported action")
	errUnsupportedFlag   error = errors.New("error, unsupported flag")
	errUnsupportedParam  error = errors.New("error, unsupported parameter")
	errWebsiteFlag       error = errors.New("error, unsupported website")
	listOfFlags          []flag
	inputParams          inputParameters
)

type inputParameters struct {
	website string
	item    string
}

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
		if parseInput(scanner.Text()) != nil {
			log.Printf("%v Failed to call parseInput()", tools.FuncName())
		}
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
			if supported, _ := checkIfFlagIsSupported(currentFlag); !supported {
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

func setParameters(paramToset, paramValue string, inputParams inputParameters) error {
	switch paramToset {
	case "website":
		inputParams.website = paramValue
	case "item":
		inputParams.item = paramValue
	default:
		log.Printf("%v Unsupported parameter %v", tools.FuncName(), paramToset)
		return errUnsupportedParam
	}
	return nil
}

func checkIfFlagIsSupported(flag string) (bool, error) {
	flag = strings.TrimSpace(flag)
	if _, flagSupported := supportedFlags[flag]; !flagSupported {
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

func callWebsiteFunction(functionToCall string, websiteToCall website.Website, params inputParameters) error {
	switch functionToCall {
	case searchAction:
		websiteToCall.SearchWebsite(params.item)
	default:
		log.Printf("%v Unsupported action %v", tools.FuncName(), functionToCall)
		return errUnsupportedAction
	}
	return nil
}

// func getWebsiteFromFlags() {
// 	log.Printf("%v Attempting to get website from flags", tools.FuncName())
// 	switch applicationWebsite {
// 	case amazonFlag:
// 		w := amazon.Constructor()
// 		//t := &amazon.Amazon{}
// 		//invokeAmazonActions(w)
// 		//invokeAmazonActions(t)
// 	case bestbuyFlag:
// 		log.Printf("Not implemented yet")
// 	default:
// 		log.Fatalf("%v Failed to complete function, website %s is not supported", tools.FuncName(), applicationWebsite)
// 	}
// }

func shutDown() {
	log.Printf("Exiting Program")
	os.Exit(1)
}
