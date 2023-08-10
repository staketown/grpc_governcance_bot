package main

import (
	"bytes"
	"cosmos_governance_bot"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const timeFormat = "Monday, January 2, 2006 03:04:05 PM (UTC)"

func main() {
	fmt.Println("Bot has been started")
	// Load configs
	cft := new(cosmos_governance_bot.Config).Load()

	for range time.Tick(time.Minute * time.Duration(cft.IntervalMinutes)) {
		runInterval(cft)
	}
}

func runInterval(cft cosmos_governance_bot.Config) {
	fmt.Println("Staring iteration at: ", time.Now().UTC().Format(timeFormat))
	allProposalsByChain := cosmos_governance_bot.LoadGovernance(&cft)

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + cft.BotToken)

	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Load proposals
	for _, proposalByChain := range allProposalsByChain {
		chainName := proposalByChain.ChainName

		for _, proposal := range proposalByChain.Proposals.Proposals {
			color, err := strconv.ParseUint(cft.Chains[chainName].Discord.HexColor, 16, 64)

			var proposalType string
			if len(proposal.Messages) == 0 {
				proposalType = "test"
			} else {
				proposalType = GetPrettyProposalType(proposal.Messages[0].TypeUrl)
			}

			var description string
			var title string
			cascadiaDetails := cosmos_governance_bot.GetIpfsData(proposal.Metadata)

			if cascadiaDetails != nil {
				description = cascadiaDetails.Details
				title = cascadiaDetails.Title
			} else {
				title = proposal.Metadata
				description = proposal.Metadata
			}

			if len(description) > 1024 {
				description = description[:1000] + " ..."
			}

			whParams := &discordgo.WebhookParams{
				Content: cft.Chains[chainName].Discord.Tags,
				AllowedMentions: &discordgo.MessageAllowedMentions{
					Parse: []discordgo.AllowedMentionType{
						discordgo.AllowedMentionTypeRoles,
					},
				},

				AvatarURL: cft.Chains[chainName].Discord.AvatarUrl,
				Embeds: []*discordgo.MessageEmbed{
					{
						Title: title,
						Color: int(color),
						Thumbnail: &discordgo.MessageEmbedThumbnail{
							URL: cft.Chains[chainName].Discord.AvatarUrl,
						},
						Fields: []*discordgo.MessageEmbedField{
							{
								Name:  "Proposal Type",
								Value: proposalType,
							},
							{
								Name:  "Proposal Details",
								Value: description,
							},
							{
								Name:  "Details",
								Value: cft.Chains[chainName].ExplorerGovUrl + strconv.FormatInt(int64(proposal.Id), 10),
							},
							{
								Name:  "Ends At",
								Value: proposal.VotingEndTime.UTC().Format(timeFormat),
							},
						},
						Footer: &discordgo.MessageEmbedFooter{
							Text: "⚡ Powered by StakeTown",
						},
					},
				},
			}

			var webHookId string
			var webHookToken string

			if !cft.Production {
				webHookId = cft.Discord.WebhookId
				webHookToken = cft.Discord.WebhookToken
			} else {
				webHookId = cft.Chains[chainName].Discord.WebhookId
				webHookToken = cft.Chains[chainName].Discord.WebhookToken
			}

			message, err := dg.WebhookExecute(webHookId, webHookToken, true, whParams)

			if err != nil {
				fmt.Errorf("can't send message, %s", err)
			} else {
				fmt.Printf("On chain '%s' has been sent new message: %s \n", chainName, message.Embeds[0].Title)
			}
		}
	}

	// Cleanly close down the Discord session.
	fmt.Println("Ending iteration at: ", time.Now().UTC().Format(timeFormat))
	dg.Close()
}

func GetPrettyProposalType(proposalType string) string {
	tmp := strings.Split(proposalType, ".")
	buf := &bytes.Buffer{}

	for i, rune := range tmp[len(tmp)-1] {
		if unicode.IsUpper(rune) && i > 0 {
			buf.WriteRune(' ')
		}
		buf.WriteRune(rune)
	}
	return buf.String()
}
