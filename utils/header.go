package utils

import "github.com/pterm/pterm"

var Version string = "undefined"

func SetVersion(version string) {
    Version = version
}

func PrintHeader() {
    // Generate BigLetters
    pterm.DefaultCenter.WithCenterEachLineSeparately().Println("")
    s, _ := pterm.DefaultBigText.WithLetters(pterm.NewLettersFromStringWithStyle("Profero", pterm.DefaultBox.TextStyle)).Srender()
    pterm.DefaultCenter.Println(s) // Print BigLetters with the default CenterPrinter
    pterm.DefaultCenter.WithCenterEachLineSeparately().Print("log4jScanner")
    pterm.DefaultCenter.WithCenterEachLineSeparately().Print(Version)
}
