package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Can/sysrow/pkg/task"
)

// Command represents a CLI command
type Command struct {
	Name        string
	Description string
	Usage       string
	Handler     func([]string)
	Subcommands map[string]*Command
	Flags       *flag.FlagSet
}

// NewCommand creates a new command
func NewCommand(name, description, usage string, handler func([]string)) *Command {
	return &Command{
		Name:        name,
		Description: description,
		Usage:       usage,
		Handler:     handler,
		Subcommands: make(map[string]*Command),
		Flags:       flag.NewFlagSet(name, flag.ExitOnError),
	}
}

// AddSubcommand adds a subcommand to the command
func (c *Command) AddSubcommand(sub *Command) {
	c.Subcommands[sub.Name] = sub
}

// Execute executes the command with the given arguments
func (c *Command) Execute(args []string) {
	if len(args) == 0 {
		c.Handler(args)
		return
	}

	subName := args[0]
	if sub, ok := c.Subcommands[subName]; ok {
		sub.Execute(args[1:])
	} else {
		c.Handler(args)
	}
}

// PrintUsage prints the usage information for the command
func (c *Command) PrintUsage() {
	fmt.Printf("%s - %s\n\n", c.Name, c.Description)
	fmt.Printf("Kullanım:\n  %s\n\n", c.Usage)

	if len(c.Subcommands) > 0 {
		fmt.Println("Alt Komutlar:")
		for name, sub := range c.Subcommands {
			fmt.Printf("  %-15s %s\n", name, sub.Description)
		}
		fmt.Println()
	}
}

func printUsage() {
	// Print app name and description
	fmt.Printf("\n%s - %s\n\n", i18n.Get("app_name"), i18n.Get("app_description"))
	
	// Print usage
	fmt.Println(i18n.Get("cli_messages.usage"))
	fmt.Print("  sysrow [command] [arguments]\n\n")

	// Print commands
	fmt.Println("Commands:")
	fmt.Printf("  %-10s %s\n", "queue", i18n.Get("commands_menu.queue"))
	fmt.Printf("  %-10s %s\n", "delay", i18n.Get("commands_menu.delay"))
	fmt.Printf("  %-10s %s\n", "run", i18n.Get("commands_menu.run"))
	fmt.Printf("  %-10s %s\n", "group", i18n.Get("commands_menu.group"))
	fmt.Printf("  %-10s %s\n", "list", i18n.Get("commands_menu.status"))
	fmt.Printf("  %-10s %s\n", "status", i18n.Get("commands_menu.status"))
	fmt.Printf("  %-10s %s\n", "logs", i18n.Get("commands_menu.status"))
	fmt.Printf("  %-10s %s\n", "cancel", i18n.Get("commands_menu.cancel"))
	fmt.Printf("  %-10s %s\n", "help", "Detailed help information")

	// Print help hint
	fmt.Printf("\n%s\n", i18n.Get("cli_messages.help_hint"))
}

func main() {
	// Initialize the application
	if err := task.InitializeDataDirectory(); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing data directory: %v\n", err)
		os.Exit(1)
	}

	// Initialize i18n system
	language := ""
	
	// Check for language flag
	for i, arg := range os.Args {
		if arg == "--lang" && i+1 < len(os.Args) {
			language = os.Args[i+1]
			// Remove the flag and its value from args
			newArgs := make([]string, 0, len(os.Args)-2)
			newArgs = append(newArgs, os.Args[:i]...)
			if i+2 < len(os.Args) {
				newArgs = append(newArgs, os.Args[i+2:]...)
			}
			os.Args = newArgs
			break
		}
	}

	// Initialize language system with a fallback to English if it fails
	if err := InitializeLanguage(language); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Could not initialize language system: %v\n", err)
		// Create a basic English translation map for critical messages
		i18n = &I18n{
			CurrentLang: "en",
			LangPath:    "./lang",
			Translations: Translations{
				"app_name":        "SysRow",
				"app_description": "CLI-based background task manager",
				"cli_messages": map[string]interface{}{
					"usage":     "Usage:",
					"help_hint": "Use 'sysrow help' for more information.",
				},
				"commands_menu": map[string]interface{}{
					"queue":  "Add a task to the queue",
					"delay":  "Schedule a task for later execution",
					"run":    "Run a task immediately, optionally in background",
					"group":  "Manage task groups",
					"status": "List tasks and check their status",
					"cancel": "Cancel a task",
				},
			},
			FallbackLang: "en",
		}
	}

	// Parse command line arguments
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(0)
	}

	// Extract the command
	cmd := strings.ToLower(os.Args[1])

	// Process the command
	switch cmd {
	case "queue":
		handleQueueCommand(os.Args[2:])
	case "delay":
		handleDelayCommand(os.Args[2:])
	case "run":
		handleRunCommand(os.Args[2:])
	case "group":
		handleGroupCommand(os.Args[2:])
	case "list":
		handleListCommand(os.Args[2:])
	case "status":
		handleStatusCommand(os.Args[2:])
	case "logs":
		handleLogsCommand(os.Args[2:])
	case "cancel":
		handleCancelCommand(os.Args[2:])
	case "help":
		showDetailedHelp()
	case "--help", "-h":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Bilinmeyen komut: %s\n\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

// Command handlers
func handleQueueCommand(args []string) {
	flags := flag.NewFlagSet("queue", flag.ExitOnError)
	priority := flags.String("priority", "normal", "Görev önceliği (low, normal, high)")
	flags.StringVar(priority, "p", "normal", "Görev önceliği (kısa form)")

	if err := flags.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Argüman ayrıştırma hatası: %v\n", err)
		os.Exit(1)
	}

	command := flags.Arg(0)
	if command == "" {
		fmt.Println("Hata: Çalıştırılacak komut belirtilmedi")
		fmt.Println("Kullanım: sysrow queue [--priority=<öncelik>] <komut>")
		os.Exit(1)
	}

	fmt.Printf("Sıraya ekleniyor: '%s' (öncelik: %s)\n", command, *priority)
	fmt.Println("Queue command not yet implemented")
}

func handleDelayCommand(args []string) {
	flags := flag.NewFlagSet("delay", flag.ExitOnError)
	at := flags.String("at", "", "Belirli bir saatte çalıştır (HH:MM formatında)")
	after := flags.String("after", "", "Belirli bir süre sonra çalıştır (5m, 2h, 1d gibi)")

	if err := flags.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Argüman ayrıştırma hatası: %v\n", err)
		os.Exit(1)
	}

	command := flags.Arg(0)
	if command == "" {
		fmt.Println("Hata: Çalıştırılacak komut belirtilmedi")
		fmt.Println("Kullanım: sysrow delay [--at=<zaman>|--after=<süre>] <komut>")
		os.Exit(1)
	}

	if *at != "" {
		fmt.Printf("Zamanlanıyor: '%s' (saat: %s)\n", command, *at)
	} else if *after != "" {
		fmt.Printf("Zamanlanıyor: '%s' (%s sonra)\n", command, *after)
	} else {
		fmt.Println("Hata: --at veya --after parametresi belirtilmedi")
		fmt.Println("Kullanım: sysrow delay [--at=<zaman>|--after=<süre>] <komut>")
		os.Exit(1)
	}

	fmt.Println("Delay command not yet implemented")
}

func handleRunCommand(args []string) {
	flags := flag.NewFlagSet("run", flag.ExitOnError)
	background := flags.Bool("bg", false, "Arka planda çalıştır")

	if err := flags.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Argüman ayrıştırma hatası: %v\n", err)
		os.Exit(1)
	}

	command := flags.Arg(0)
	if command == "" {
		fmt.Println("Hata: Çalıştırılacak komut belirtilmedi")
		fmt.Println("Kullanım: sysrow run [--bg] <komut>")
		os.Exit(1)
	}

	if *background {
		fmt.Printf("Arka planda çalıştırılıyor: '%s'\n", command)
	} else {
		fmt.Printf("Çalıştırılıyor: '%s'\n", command)
	}

	fmt.Println("Run command not yet implemented")
}

func handleGroupCommand(args []string) {
	if len(args) == 0 {
		fmt.Println("Hata: Alt komut belirtilmedi")
		fmt.Println("Kullanım: sysrow group <create|add|run|delete> [argümanlar]")
		os.Exit(1)
	}

	subCmd := args[0]
	subArgs := args[1:]

	switch subCmd {
	case "create":
		handleGroupCreateCommand(subArgs)
	case "add":
		handleGroupAddCommand(subArgs)
	case "run":
		handleGroupRunCommand(subArgs)
	case "delete":
		handleGroupDeleteCommand(subArgs)
	default:
		fmt.Printf("Bilinmeyen alt komut: %s\n", subCmd)
		fmt.Println("Kullanım: sysrow group <create|add|run|delete> [argümanlar]")
		os.Exit(1)
	}
}

func handleGroupCreateCommand(args []string) {
	if len(args) == 0 {
		fmt.Println("Hata: Grup adı belirtilmedi")
		fmt.Println("Kullanım: sysrow group create <grup_adı>")
		os.Exit(1)
	}

	groupName := args[0]
	fmt.Printf("Grup oluşturuluyor: '%s'\n", groupName)
	fmt.Println("Group create command not yet implemented")
}

func handleGroupAddCommand(args []string) {
	if len(args) < 2 {
		fmt.Println("Hata: Grup adı veya komut belirtilmedi")
		fmt.Println("Kullanım: sysrow group add <grup_adı> <komut>")
		os.Exit(1)
	}

	groupName := args[0]
	command := args[1]
	fmt.Printf("Gruba ekleniyor: '%s' -> '%s'\n", command, groupName)
	fmt.Println("Group add command not yet implemented")
}

func handleGroupRunCommand(args []string) {
	if len(args) == 0 {
		fmt.Println("Hata: Grup adı belirtilmedi")
		fmt.Println("Kullanım: sysrow group run <grup_adı>")
		os.Exit(1)
	}

	groupName := args[0]
	fmt.Printf("Grup çalıştırılıyor: '%s'\n", groupName)
	fmt.Println("Group run command not yet implemented")
}

func handleGroupDeleteCommand(args []string) {
	if len(args) == 0 {
		fmt.Println("Hata: Grup adı belirtilmedi")
		fmt.Println("Kullanım: sysrow group delete <grup_adı>")
		os.Exit(1)
	}

	groupName := args[0]
	fmt.Printf("Grup siliniyor: '%s'\n", groupName)
	fmt.Println("Group delete command not yet implemented")
}

func handleListCommand(args []string) {
	fmt.Println("Görevler listeleniyor...")
	fmt.Println("List command not yet implemented")
}

func handleStatusCommand(args []string) {
	if len(args) == 0 {
		fmt.Println("Hata: Görev ID'si belirtilmedi")
		fmt.Println("Kullanım: sysrow status <görev_id>")
		os.Exit(1)
	}

	taskID := args[0]
	fmt.Printf("Görev durumu kontrol ediliyor: '%s'\n", taskID)
	fmt.Println("Status command not yet implemented")
}

func handleLogsCommand(args []string) {
	if len(args) == 0 {
		fmt.Println("Hata: Görev ID'si belirtilmedi")
		fmt.Println("Kullanım: sysrow logs <görev_id>")
		os.Exit(1)
	}

	taskID := args[0]
	fmt.Printf("Görev günlükleri görüntüleniyor: '%s'\n", taskID)
	fmt.Println("Logs command not yet implemented")
}

func handleCancelCommand(args []string) {
	if len(args) == 0 {
		fmt.Println("Hata: Görev ID'si belirtilmedi")
		fmt.Println("Kullanım: sysrow cancel <görev_id>")
		os.Exit(1)
	}

	taskID := args[0]
	fmt.Printf("Görev iptal ediliyor: '%s'\n", taskID)
	fmt.Println("Cancel command not yet implemented")
}
