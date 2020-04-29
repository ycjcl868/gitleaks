package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"regexp"
	"testing"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

type TomlConfig struct {
	Title string
	Whitelist struct {
		Description string
		Commits     []string
		Files       []string
		Paths       []string
	}
	Rules []struct {
		Description   string
		Regex         string
		FileNameRegex string
		FilePathRegex string
		Tags          []string
		Entropies     []struct {
			Min   string
			Max   string
			Group string
		}
		Whitelist []struct {
			Description string
			Regex       string
			File        string
			Path        string
		}
	}
}

func getData() (tomlconfig TomlConfig, jsonData map[string]map[string]bool) {

	var cwd, _ = os.Getwd()
	var config = TomlConfig{}
	var _, err = toml.DecodeFile(path.Join(cwd, "../.gitleaks.toml"), &config)
	if err != nil {
		fmt.Println(err)
	}
	var jsonFile, errJSONFile = os.Open(path.Join(cwd, "../package-lock.json"))
	if errJSONFile != nil {
		fmt.Println(errJSONFile)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result map[string]map[string]bool
	json.Unmarshal([]byte(byteValue), &result)
	return config, result
}


func TestHello(t *testing.T)  {
	config, jsonData := getData()

	t.Log(config.Rules[0].Description)
	for _, rule := range config.Rules {
		ruleCase := jsonData[rule.Description]
		ruleExp := rule.Regex
		if ruleCase != nil {
			//t.Log(ruleCase)
			//t.Log(ruleExp)
			for testString := range ruleCase {
				expectVal := ruleCase[testString]

				//t.Log(ruleCase)
				re, err := regexp.Compile(ruleExp)
				if err != nil {
					fmt.Println(err)
				}
				match := re.FindString(testString)
				t.Log("===start==")
				t.Log("testString", testString)
				t.Log("expectVal", expectVal)
				t.Log("match", match)
				actualVal := match != ""

				t.Log("actualVal", actualVal)
				t.Log("===end==")
				if (expectVal != actualVal) {
					t.Errorf("<%v> test failed, case Value %v, expected <%v>, but get <%v>", rule.Description, testString, expectVal, actualVal)
				}
			}
		}
	}
}
