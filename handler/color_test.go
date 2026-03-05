package handler

import (
	"fmt"
	"testing"
)

func TestColor(t *testing.T) {
	fmt.Print(Bold + "bold" + Reset + "\n")
	fmt.Print(Faint + "faint" + Reset + "\n")
	fmt.Print(Underlined + "underlined" + Reset + "\n")
	fmt.Print(Blink + "blink" + Reset + "\n")

	fmt.Print(Black + "black" + Reset + "\n")
	fmt.Print(Red + "red" + Reset + "\n")
	fmt.Print(Green + "green" + Reset + "\n")
	fmt.Print(Yellow + "yellow" + Reset + "\n")
	fmt.Print(Blue + "blue" + Reset + "\n")
	fmt.Print(Purple + "purple" + Reset + "\n")
	fmt.Print(Cyan + "cyan" + Reset + "\n")
	fmt.Print(Gray + "gray" + Reset + "\n")

	fmt.Print(White + "white logging" + Reset + "\n")

	fmt.Print(GrayBlack + "gray background with black foreground" + Reset + "\n")
	fmt.Print(CyanRed + "cyan background with red foreground" + Reset + "\n")
	fmt.Print(PurpleGreen + "purple background with green foreground" + Reset + "\n")
	fmt.Print(BlueYellow + "blue background with yellow foreground" + Reset + "\n")
	fmt.Print(YellowBlue + "yellow background with blue foreground" + Reset + "\n")
	fmt.Print(GreenPurple + "green background with purple foreground" + Reset + "\n")
	fmt.Print(RedCyan + "red background with cyan foreground" + Reset + "\n")
	fmt.Print(BlackGray + "black background with gray foreground" + Reset + "\n")

	fmt.Print(BlackWhite + "black background with white foreground" + Reset + "\n")

}
