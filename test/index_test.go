package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"testing"

	"github.com/BurntSushi/toml"
)

type TomlConfig struct {
	Title     string
	Whitelist struct {
		Description string
		Commits     []string
		Files       []string
		File        string
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

func expectRule(reg string, testString string, expectBool bool, t *testing.T) {
	//t.Log(ruleCase)
	re, err := regexp.Compile(reg)
	if err != nil {
		fmt.Println(err)
	}
	actualBool := re.FindString(testString) != ""
	if actualBool != expectBool {
		t.Errorf("<%v> test failed, case Value %v, expected <%v>, but get <%v>", reg, testString, expectBool, actualBool)
	}
}


func TestRules(t *testing.T) {
	config, jsonData := getData()

	for _, rule := range config.Rules {
		ruleCase := jsonData[rule.Description]
		ruleExp := rule.Regex
		ruleFileNameExp := rule.FileNameRegex
		if ruleCase != nil {
			for testString := range ruleCase {
				expectVal := ruleCase[testString]
				if ruleFileNameExp != "" {
					expectRule(ruleFileNameExp, testString, expectVal, t)
				}
				if ruleExp != "" {
					expectRule(ruleExp, testString, expectVal, t)
				}
			}
		}
	}
}

func TestWhiteList(t *testing.T) {
	config, jsonData := getData()

	ruleCase := jsonData["WhiteList"]
	ruleExps := config.Whitelist.Files
	if ruleCase != nil {
		match := false
		for testString := range ruleCase {
			expectVal := ruleCase[testString]
			for _, ruleExp := range ruleExps {
				//t.Log(ruleCase)
				re, err := regexp.Compile(ruleExp)
				if err != nil {
					fmt.Println(err)
				}
				if match == false {
					match = re.FindString(testString) != ""
				}
			}
			if match != expectVal {
				t.Errorf("whitelist test failed, case Value %v, expected <%v>, but get <%v>", testString, expectVal, match)
			}
		}
	}
}
