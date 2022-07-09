package oid_tests

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
	"testing"

	"gopkg.in/yaml.v2"
)

// Struct to hold the parsed yaml data
type TestData struct {
	Oids []string `yaml:"trap-type-oid-prefix"`
}

var pathToCLI = "../bin/oid.sh"

func Test_YamlFileWithValidOids(t *testing.T) {
	var testData TestData
	loadYaml("../snmp.yaml", testData)
	// Execute the shell script under test
	cmdProc := exec.Command("sh", pathToCLI)

	// TODO: Could be quite slow with larger data sets
	for _, oid := range testData.Oids {
		t.Run(fmt.Sprintf("Testing OID %s", oid), func(t *testing.T) {
			output, _ := checkOid(cmdProc, oid)

			// TODO: Might be better to use the CLI output in the error output
			if strings.Contains(string(output), fmt.Sprintf("%s: false", oid)) {
				t.Errorf(fmt.Sprintf("OID %s returned 'false'. Should have returned 'true'.", oid))
			}
		})
	}

	quit(cmdProc)
}

func Test_YamlFileWithInvaidOids(t *testing.T) {
	cmdProc := exec.Command("sh", pathToCLI, "-f", "../invalid-snmp.yaml")
	output, _ := cmdProc.CombinedOutput()
	expected := "Error: Loader: no valid oids found under key 'trap-type-oid-prefix'"

	if !strings.Contains(string(output), expected) {
		t.Errorf(fmt.Sprintf("Unexpected output. Should display: %s", expected))
	}
}

func Test_NonExistentYamlFile(t *testing.T) {
	cmdProc := exec.Command("sh", pathToCLI, "-f", "./nonexistent.yaml")
	output, _ := cmdProc.CombinedOutput()
	expected := "Error: Loader: cannot process file './nonexistent.yaml'."

	if !strings.Contains(string(output), expected) {
		t.Errorf(fmt.Sprintf("Unexpected output. Should display: %s", expected))
	}
}

func Test_InvalidInput(t *testing.T) {
	cmdProc := exec.Command("sh", pathToCLI)
	output, _ := checkOid(cmdProc, "foobar")
	expected := "Error: OID 'foobar' is invalid."

	if !strings.Contains(string(output), expected) {
		t.Errorf(fmt.Sprintf("Unexpected output. Should display: %s", expected))
	}

	quit(cmdProc)
}

func loadYaml(pathToYaml string, testData TestData) {
	data, err := ioutil.ReadFile(pathToYaml)

	if err != nil {
		log.Fatal(err)
	}

	if err := testData.ParseYaml(data); err != nil {
		log.Fatal(err)
	}
}

func checkOid(cmdProc *exec.Cmd, oid string) ([]byte, error) {
	cmdProc.Stdin = strings.NewReader(fmt.Sprintf("%s\n", oid))
	output, error := cmdProc.CombinedOutput()
	return output, error
}

func quit(cmdProc *exec.Cmd) {
	cmdProc.Stdin = strings.NewReader("quit\n")
}

func (c *TestData) ParseYaml(data []byte) error {
	if err := yaml.Unmarshal(data, c); err != nil {
		return err
	}

	if len(c.Oids) == 0 {
		return errors.New("Parse Error: List of OIDs was empty")
	}

	return nil
}
