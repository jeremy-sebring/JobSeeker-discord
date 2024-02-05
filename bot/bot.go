package Bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	JobHunter "sebring.dev/JobSeeker-discord/JobHunter/v2"

	"github.com/bwmarrin/discordgo"
)

func checkNilErr(e error) {
	if e != nil {
		log.Fatal(e.Error())
	}
}

func ChunkS(s string, chunkSize int) []string {
	if len(s) == 0 {
		return nil
	}
	if chunkSize >= len(s) {
		return []string{s}
	}
	var chunks []string = make([]string, 0, (len(s)-1)/chunkSize+1)
	currentLen := 0
	currentStart := 0
	for i := range s {
		if currentLen == chunkSize {
			chunks = append(chunks, s[currentStart:i])
			currentLen = 0
			currentStart = i
		}
		currentLen++
	}
	chunks = append(chunks, s[currentStart:])
	return chunks
}

func SendSafeMessage(discord *discordgo.Session, channelID string, message string) {
	if len(message) > 2000 {
		var msgDescs = ChunkS(message, 2000)

		for d := 0; d < len(msgDescs); d++ {
			msgDesc, errDesc := discord.ChannelMessageSend(channelID, msgDescs[d])
			checkNilErr(errDesc)
			log.Println(msgDesc.ID)
		}
	} else {
		msgDesc, errDesc := discord.ChannelMessageSend(channelID, message)
		checkNilErr(errDesc)
		log.Println(msgDesc.ID)
	}
}

func Run() {
	BotToken, TokenExists := os.LookupEnv("DISCORD_TOKEN")

	if !TokenExists {
		log.Fatal("BotToken is not present")
	}

	// create a session
	discord, err := discordgo.New("Bot " + BotToken)
	checkNilErr(err)

	discord.Open()

	discord.AddHandler(newMessage)

	defer discord.Close() // close session, after function termination

	// keep bot running until there is NO os interruption (ctrl + C)
	fmt.Println("Bot running....")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

}

func newMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {

	/* prevent bot responding to its own message
	this is achieved by looking into the message author id
	if message.author.id is same as bot.author.id then just return
	*/
	if message.Author.ID == discord.State.User.ID {
		return
	}

	log.Printf(message.ChannelID)
	log.Printf(message.ID)
	// respond to user message if it contains `!help` or `!bye`
	switch {
	case message.ChannelID != "1203005434252890214":
		log.Printf("Not reacting to message")
	}

}

func CreateJobthreads(jobs []JobHunter.JobsResult) {
	const channelId = "1200466170399166555"

	const prodId = "1200466170399166555"
	log.Println(prodId)

	BotToken, TokenExists := os.LookupEnv("DISCORD_TOKEN")

	if !TokenExists {
		log.Fatal("BotToken is not present")
	}

	discord, err := discordgo.New("Bot " + BotToken)
	checkNilErr(err)
	discord.Open()
	defer discord.Close()

	for i := 0; i < len(jobs); i++ {

		var job = jobs[i]

		var threadTitle = fmt.Sprintf("%s at %s", job.Title, job.CompanyName)

		thrd, err := discord.ThreadStart(
			channelId,
			threadTitle,
			discordgo.ChannelType(11),
			1440,
		)
		checkNilErr(err)

		var location = fmt.Sprintf("Location: %s", job.Location)

		SendSafeMessage(discord, thrd.ID, location)

		var Via = fmt.Sprintf("Via: %s", job.Via)

		SendSafeMessage(discord, thrd.ID, Via)

		var jobDes = fmt.Sprintf("## Job Description: \n%s", job.Description)

		SendSafeMessage(discord, thrd.ID, jobDes)

		for x := 0; x < len(job.JobHighlights); x++ {
			var highlight = job.JobHighlights[x]
			var highlightMsg = fmt.Sprintf("## %s \n", highlight.Title)

			for y := 0; y < len(highlight.Items); y++ {
				highlightMsg = fmt.Sprintf("%s- %s\n", highlightMsg, highlight.Items[y])
			}
			SendSafeMessage(discord, thrd.ID, highlightMsg)
		}

		var linkList string = "## Links:\n"

		for l := 0; l < len(job.RelatedLinks); l++ {
			var link JobHunter.RelatedLink = job.RelatedLinks[l]
			linkList = fmt.Sprintf("%s- [%s](%s)\n", linkList, link.Text, link.Link)
		}

		SendSafeMessage(discord, thrd.ID, linkList)
	}
}
