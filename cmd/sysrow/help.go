package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Global i18n instance
var i18n *I18n

// createBox creates a box around text
func createBox(text string, width int) string {
	padding := width - len(text) - 2
	leftPad := padding / 2
	rightPad := padding - leftPad

	topBottom := "╔" + strings.Repeat("═", width-2) + "╗"
	middle := "║" + strings.Repeat(" ", leftPad) + text + strings.Repeat(" ", rightPad) + "║"
	bottom := "╚" + strings.Repeat("═", width-2) + "╝"

	return topBottom + "\n" + middle + "\n" + bottom
}



// showDetailedHelp displays comprehensive help information about the system
func showDetailedHelp() {
	// i18n should already be initialized in main.go
	// Just show the help menu
	showHelpMenu()
}

// showHelpMenu displays the main help menu and handles user interaction
func showHelpMenu() {
	// Create the header box
	headerBox := createBox(i18n.Get("help_title"), 60)

	fmt.Println("\n" + headerBox)
	fmt.Println("\n" + i18n.Get("app_description"))

	for {
		fmt.Println("\n" + i18n.Get("main_menu.title"))
		fmt.Println("1 - " + i18n.Get("main_menu.basic_commands"))
		fmt.Println("2 - " + i18n.Get("main_menu.example_scenarios"))
		fmt.Println("3 - " + i18n.Get("main_menu.additional_info"))
		fmt.Println("l - " + "Change Language / Dil Değiştir")
		fmt.Println("q - " + i18n.Get("main_menu.exit"))
		fmt.Print("\n" + i18n.Get("main_menu.prompt"))

		reader := bufio.NewReader(os.Stdin)
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			showBasicCommands()
		case "2":
			showExampleScenarios()
		case "3":
			showAdditionalInfo()
		case "l", "L":
			changeLanguage()
		case "q", "Q":
			return
		default:
			fmt.Println("\n" + i18n.Get("main_menu.invalid_choice"))
		}
	}
}

// changeLanguage allows the user to change the current language
func changeLanguage() {
	// Get available languages
	languages, err := i18n.AvailableLanguages()
	if err != nil {
		fmt.Printf("Error: Could not get available languages: %v\n", err)
		return
	}

	// Display language options
	fmt.Println("\nAvailable Languages / Mevcut Diller:")
	for i, lang := range languages {
		fmt.Printf("%d - %s\n", i+1, lang)
	}
	fmt.Println("b - Back / Geri")
	fmt.Print("\nSelect a language / Bir dil seçin: ")

	// Get user choice
	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	// Handle user choice
	if choice == "b" || choice == "B" {
		return
	}

	// Try to parse the choice as a number
	var index int
	_, err = fmt.Sscanf(choice, "%d", &index)
	if err != nil || index < 1 || index > len(languages) {
		fmt.Println("Invalid choice / Geçersiz seçim")
		return
	}

	// Load the selected language
	lang := languages[index-1]
	err = i18n.LoadLanguage(lang)
	if err != nil {
		fmt.Printf("Error: Could not load language '%s': %v\n", lang, err)
		return
	}

	fmt.Printf("Language changed to '%s' / Dil '%s' olarak değiştirildi\n", lang, lang)
}

// showBasicCommands displays information about basic commands
func showBasicCommands() {
	for {
		// Create the header box
		headerBox := createBox(i18n.Get("commands_menu.title"), 60)

		fmt.Println("\n" + headerBox)
		fmt.Println("\n" + i18n.Get("commands_menu.prompt"))
		fmt.Println("\n1 - " + i18n.Get("commands_menu.queue"))
		fmt.Println("2 - " + i18n.Get("commands_menu.delay"))
		fmt.Println("3 - " + i18n.Get("commands_menu.group"))
		fmt.Println("4 - " + i18n.Get("commands_menu.run"))
		fmt.Println("5 - " + i18n.Get("commands_menu.status"))
		fmt.Println("6 - " + i18n.Get("commands_menu.cancel"))
		fmt.Println("\nb - " + i18n.Get("navigation.back"))
		fmt.Println("q - " + i18n.Get("navigation.exit"))

		fmt.Print("\n" + i18n.Get("main_menu.prompt"))
		reader := bufio.NewReader(os.Stdin)
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			showQueueCommand()
		case "2":
			showDelayCommand()
		case "3":
			showGroupCommand()
		case "4":
			showRunCommand()
		case "5":
			showStatusCommands()
		case "6":
			showCancelCommand()
		case "b", "B":
			return
		case "q", "Q":
			os.Exit(0)
		default:
			fmt.Println("\n" + i18n.Get("main_menu.invalid_choice"))
		}
	}
}

// showQueueCommand displays information about the queue command
func showQueueCommand() {
	fmt.Println("\n" + i18n.Get("command_details.queue.title"))
	fmt.Println("   " + i18n.Get("command_details.queue.description"))
	fmt.Println("   ")
	fmt.Println("   " + i18n.Get("command_details.queue.usage"))
	fmt.Println("   " + i18n.Get("command_details.queue.example1"))
	fmt.Println("   " + i18n.Get("command_details.queue.example2"))
	fmt.Println("   ")
	fmt.Println("   " + i18n.Get("command_details.queue.options"))
	fmt.Println("   " + i18n.Get("command_details.queue.option_priority"))

	fmt.Println("\n" + i18n.Get("navigation.continue"))
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

// showDelayCommand displays information about the delay command
func showDelayCommand() {
	fmt.Println("\n" + i18n.Get("command_details.delay.title"))
	fmt.Println("   " + i18n.Get("command_details.delay.description"))
	fmt.Println("   ")
	fmt.Println("   " + i18n.Get("command_details.delay.usage"))
	fmt.Println("   " + i18n.Get("command_details.delay.example1"))
	fmt.Println("   " + i18n.Get("command_details.delay.example2"))
	fmt.Println("   ")
	fmt.Println("   " + i18n.Get("command_details.delay.options"))
	fmt.Println("   " + i18n.Get("command_details.delay.option_at"))
	fmt.Println("   " + i18n.Get("command_details.delay.option_after"))

	fmt.Println("\n" + i18n.Get("navigation.continue"))
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

// showGroupCommand displays information about the group command
func showGroupCommand() {
	fmt.Println("\n" + i18n.Get("command_details.group.title"))
	fmt.Println("   " + i18n.Get("command_details.group.description"))
	fmt.Println("   ")
	fmt.Println("   " + i18n.Get("command_details.group.usage"))
	fmt.Println("   " + i18n.Get("command_details.group.example1"))
	fmt.Println("   " + i18n.Get("command_details.group.example2"))
	fmt.Println("   " + i18n.Get("command_details.group.example3"))
	fmt.Println("   " + i18n.Get("command_details.group.example4"))
	fmt.Println("   " + i18n.Get("command_details.group.example5"))

	fmt.Println("\n" + i18n.Get("navigation.continue"))
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

// showRunCommand displays information about the run command
func showRunCommand() {
	fmt.Println("\n" + i18n.Get("command_details.run.title"))
	fmt.Println("   " + i18n.Get("command_details.run.description"))
	fmt.Println("   ")
	fmt.Println("   " + i18n.Get("command_details.run.usage"))
	fmt.Println("   " + i18n.Get("command_details.run.example1"))

	fmt.Println("\n" + i18n.Get("navigation.continue"))
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

// showStatusCommands displays information about the list, status, and logs commands
func showStatusCommands() {
	fmt.Println("\n" + i18n.Get("command_details.status.title"))
	fmt.Println("   " + i18n.Get("command_details.status.description"))
	fmt.Println("   ")
	fmt.Println("   " + i18n.Get("command_details.status.usage"))
	fmt.Println("   " + i18n.Get("command_details.status.example1"))
	fmt.Println("   " + i18n.Get("command_details.status.example2"))
	fmt.Println("   " + i18n.Get("command_details.status.example3"))

	fmt.Println("\n" + i18n.Get("navigation.continue"))
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

// showCancelCommand displays information about the cancel command
func showCancelCommand() {
	fmt.Println("\n" + i18n.Get("command_details.cancel.title"))
	fmt.Println("   " + i18n.Get("command_details.cancel.description"))
	fmt.Println("   ")
	fmt.Println("   " + i18n.Get("command_details.cancel.usage"))
	fmt.Println("   " + i18n.Get("command_details.cancel.example1"))

	fmt.Println("\n" + i18n.Get("navigation.continue"))
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

// showExampleScenarios displays example scenarios
func showExampleScenarios() {
	// Create the header box
	headerBox := createBox(i18n.Get("examples.title"), 60)

	fmt.Println("\n" + headerBox)

	fmt.Println("\n1. " + i18n.Get("examples.backup.title"))
	fmt.Println("   " + i18n.Get("examples.backup.command"))

	fmt.Println("\n2. " + i18n.Get("examples.build.title"))
	fmt.Println("   " + i18n.Get("examples.build.command"))

	fmt.Println("\n3. " + i18n.Get("examples.maintenance.title"))
	fmt.Println("   " + i18n.Get("examples.maintenance.command1"))
	fmt.Println("   " + i18n.Get("examples.maintenance.command2"))
	fmt.Println("   " + i18n.Get("examples.maintenance.command3"))
	fmt.Println("   " + i18n.Get("examples.maintenance.command4"))
	fmt.Println("   " + i18n.Get("examples.maintenance.command5"))

	fmt.Println("\n4. " + i18n.Get("examples.background.title"))
	fmt.Println("   " + i18n.Get("examples.background.command"))

	fmt.Println("\n5. " + i18n.Get("examples.reboot.title"))
	fmt.Println("   " + i18n.Get("examples.reboot.command"))

	fmt.Println("\n" + i18n.Get("navigation.return_to_main"))
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

// showAdditionalInfo displays additional information
func showAdditionalInfo() {
	// Create the header box
	headerBox := createBox(i18n.Get("additional_info.title"), 60)

	fmt.Println("\n" + headerBox)

	fmt.Println("\n- " + i18n.Get("additional_info.point1"))
	fmt.Println("- " + i18n.Get("additional_info.point2"))
	fmt.Println("- " + i18n.Get("additional_info.point3"))
	fmt.Println("- " + i18n.Get("additional_info.point4"))
	fmt.Println("- " + i18n.Get("additional_info.point5"))

	fmt.Println("\n" + i18n.Get("additional_info.more_info"))

	fmt.Println("\n" + i18n.Get("navigation.return_to_main"))
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
