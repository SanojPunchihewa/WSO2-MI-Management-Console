/*
*  Copyright (c) WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
*
*  WSO2 Inc. licenses this file to you under the Apache License,
*  Version 2.0 (the "License"); you may not use this file except
*  in compliance with the License.
*  You may obtain a copy of the License at
*
*    http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing,
* software distributed under the License is distributed on an
* "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
* KIND, either express or implied.  See the License for the
* specific language governing permissions and limitations
* under the License.
*/

package cmd

import (
	"errors"
    "mi-management-cli/utils"
    "github.com/spf13/cobra"
    "github.com/renstrom/dedent"
    "net/http"
	"encoding/xml"
	"os"
	"github.com/olekukonko/tablewriter"
)

var sequenceName string

// Show Sequence command related usage info
const showSequenceCmdLiteral = "sequence"
const showSequenceCmdShortDesc = "Get information about the specified Sequence"

var showSequenceCmdLongDesc = "Get information about the Sequence specified by the flag --name, -n\n"

var showSequenceCmdExamples = dedent.Dedent(`
Example:
  ` + utils.ProjectName + ` ` + showCmdLiteral + ` ` + showSequenceCmdLiteral + ` -n TestSequence
`)

// sequenceShowCmd represents the show sequence command
var sequenceShowCmd = &cobra.Command{
	Use:   showSequenceCmdLiteral,
	Short: showSequenceCmdShortDesc,
	Long: showSequenceCmdLongDesc + showSequenceCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo+"Show sequence called")
		executeGetSequenceCmd(sequenceName)
	},
}

func init() {
	showCmd.AddCommand(sequenceShowCmd)

	// Here you will define your flags and configuration settings.

	sequenceShowCmd.Flags().StringVarP(&sequenceName, "name", "n", "", "Name of the Sequence")
    sequenceShowCmd.MarkFlagRequired("name")
}

func executeGetSequenceCmd(sequencename string) {

    sequence, err := GetSequenceInfo(sequencename)

    if err == nil {
        // Printing the details of the Sequence
        printSequenceInfo(*sequence)
        
    } else {
        utils.Logln(utils.LogPrefixError+"Getting Information of the Sequence", err)
    }

    // if flagExportAPICmdToken != "" {
    //  // token provided with --token (-t) flag
    //  if exportAPICmdUsername != "" || exportAPICmdPassword != "" {
    //      // username and/or password provided with -u and/or -p flags
    //      // Error
    //      utils.HandleErrorAndExit("username/password provided with OAuth token.", nil)
    //  } else {
    //      // token only, proceed with token
    //  }
    // } else {
    //  // no token provided with --token (-t) flag
    //  // proceed with username and password
    //  accessToken, apiManagerEndpoint, preCommandErr := utils.ExecutePreCommand(listApisCmdEnvironment, listApisCmdUsername,
    //      listApisCmdPassword, utils.MainConfigFilePath, utils.EnvKeysAllFilePath)

    //  if preCommandErr == nil {
    //      if listApisCmdQuery != "" {
    //          fmt.Println("Search query:", listApisCmdQuery)
    //      }
    //      count, apis, err := GetCarbonAppInfo(listApisCmdQuery, accessToken, apiManagerEndpoint)

    //      if err == nil {
    //          // Printing the list of available APIs
    //          fmt.Println("Environment:", listApisCmdEnvironment)
    //          fmt.Println("No. of APIs:", count)
    //          if count > 0 {
    //              printAPIs(apis)
    //          }
    //      } else {
    //          utils.Logln(utils.LogPrefixError+"Getting List of APIs", err)
    //      }
    //  } else {
    //      utils.HandleErrorAndExit("Error calling '"+listCmdLiteral+" "+apisCmdLiteral+"'", preCommandErr)
    //  }
    // }
}

// GetSequenceInfo
// @param name of the sequence
// @return Sequence Object
// @return error
func GetSequenceInfo(name string) (*utils.Sequence, error) {

    finalUrl := utils.RESTAPIBase + utils.PrefixSequences + "?inboundEndpointName=" + name

    utils.Logln(utils.LogPrefixInfo+"URL:", finalUrl)

    headers := make(map[string]string)
    // headers[utils.HeaderAuthorization] = utils.HeaderValueAuthPrefixBearer + " " + accessToken

    resp, err := utils.InvokeGETRequest(finalUrl, headers)

    if err != nil {
        utils.HandleErrorAndExit("Unable to connect to "+finalUrl, err)
    }

    utils.Logln(utils.LogPrefixInfo+"Response:", resp.Status())

    if resp.StatusCode() == http.StatusOK {
        sequenceResponse := &utils.Sequence{}
        unmarshalError := xml.Unmarshal([]byte(resp.Body()), &sequenceResponse)

        if unmarshalError != nil {
            utils.HandleErrorAndExit(utils.LogPrefixError+"invalid XML response", unmarshalError)
        }

        return sequenceResponse, nil
    } else {
        return nil, errors.New(resp.Status())
    }
}

// printSequenceInfo
// @param task : Sequence object
func printSequenceInfo(sequence utils.Sequence) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	data := []string{"NAME", sequence.Name}
	table.Append(data)

	data = []string{"CONTAINER", sequence.Container}
    table.Append(data)
    
    for _, mediator := range sequence.Mediators {
        data = []string{"MEDIATORS", mediator}
		table.Append(data)
	}

	table.SetBorders(tablewriter.Border{Left: true, Top: true, Right: true, Bottom: false})
    table.SetRowLine(true) 
    table.SetAutoMergeCells(true)
	table.Render() // Send output
}