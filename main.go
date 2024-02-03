package main

// client id = 1203330156320653363
import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo" // Import DiscordGo library
)

var (
	token string // Your bot token
)

func init() {
	//bot token
	token = "MTIwMzMzMDE1NjMyMDY1MzM2Mw.GUVFGJ.4cow_ViD2ZYCfwsxzJ8Un7rehvc2y2RtBrBYTg"
}

func main() {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token) // Create new Discord session
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate) // Register messageCreate function to handle messages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open() // Open Discord connection
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
		return
	}

	fmt.Println("Bot is now running. Press CTRL+C to exit.")

	// Wait here until CTRL-C or other term signal is received.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close() // Close Discord session
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages from the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Parse the command
	if strings.HasPrefix(m.Content, "!") {
		args := strings.Fields(m.Content[1:])
		command := strings.ToLower(args[0])
		switch command {
		case "help":
			helpCommand(s, m.ChannelID)
		case "poll":
			pollCommand(s, m.ChannelID, args[1:])
		default:
			s.ChannelMessageSend(m.ChannelID, "Unknown command. Type !help for list of commands.")
		}
	}
}

func helpCommand(s *discordgo.Session, channelID string) {
	helpMessage := `
**Available Commands:**
- !help: Displays this help message.
- !poll <question>: Creates a poll with the provided question and yes/no options.
`
	s.ChannelMessageSend(channelID, helpMessage)
}

func pollCommand(s *discordgo.Session, channelID string, args []string) {
	if len(args) < 1 {
		s.ChannelMessageSend(channelID, "Usage: !poll <question>")
		return
	}

	question := strings.Join(args, " ")
	pollMessage := fmt.Sprintf("**Poll:** %s\nReact with :thumbsup: or :thumbsdown: to vote.", question)
	msg, err := s.ChannelMessageSend(channelID, pollMessage)
	if err != nil {
		fmt.Println("Error sending poll message: ", err)
		return
	}

	// Add reactions for voting
	err = s.MessageReactionAdd(channelID, msg.ID, "👍")
	if err != nil {
		fmt.Println("Error adding thumbs up reaction: ", err)
	}
	err = s.MessageReactionAdd(channelID, msg.ID, "👎")
	if err != nil {
		fmt.Println("Error adding thumbs down reaction: ", err)
	}
}
