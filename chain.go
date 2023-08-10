package cosmos_governance_bot

import (
	"encoding/json"
	"fmt"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"os"
	"sort"
)

const FileName = "chains.json"

func LoadGovernance(cft *Config) []Proposals {
	data, err := os.ReadFile(FileName)
	if err != nil {
		fmt.Println(err)
	}

	var chains map[string]int64
	err = json.Unmarshal(data, &chains)
	if err != nil {
		fmt.Println(err)
	}

	var allProposalsByChainName []Proposals

	for chainName, chain := range cft.Chains {
		fmt.Printf("Getting governance from chain: %s\n", chainName)

		proposalsResponse, err := GetGrpcMessage(chain.GrpcUrl, chains[chainName])

		if err != nil {
			fmt.Println(err)
			break
		}

		sort.Slice(proposalsResponse.Proposals, func(i, j int) bool {
			return proposalsResponse.Proposals[i].Id > proposalsResponse.Proposals[j].Id
		})

		if len(proposalsResponse.Proposals) > 0 {
			latestProposalId := proposalsResponse.Proposals[0].Id
			if err != nil {
				fmt.Println(err)
			}
			chains[chainName] = int64(latestProposalId)

			proposal := Proposals{
				ChainName: chainName,
				Proposals: proposalsResponse,
			}

			allProposalsByChainName = append(allProposalsByChainName, proposal)
		}
	}

	updatedData, err := json.Marshal(chains)

	if err != nil {
		fmt.Println(err)
	}

	err = os.WriteFile(FileName, updatedData, 0644)

	if err != nil {
		fmt.Println(err)
	}

	return allProposalsByChainName
}

type Proposals struct {
	Proposals *v1.QueryProposalsResponse
	ChainName string
}
