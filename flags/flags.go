package flags

import (
	"github.com/robjporter/go-functions/as"
	"github.com/robjporter/go-functions/kingpin"
)

var (
	add    = kingpin.Command("add", "Register a new UCS domain.")
	update = kingpin.Command("update", "Update a UCS domain.")
	delete = kingpin.Command("delete", "Remove a UCS domain.")
	show   = kingpin.Command("show", "Show a UCS domain.")
	run    = kingpin.Command("run", "Run the main application.")
	output = kingpin.Command("output", "Configure the output file.")
	input  = kingpin.Command("input", "Configure the input file.")

	debug = kingpin.Flag("debug", "Enable debug mode.").Bool()

	addUCS    = add.Command("ucs", "Add a UCS Domain")
	updateUCS = update.Command("ucs", "Update a UCS Domain")
	deleteUCS = delete.Command("ucs", "Delete a UCS Domain")
	showUCS   = show.Command("ucs", "Show a UCS Domain")

	showAll = show.Command("all", "Show all")

	addUCSPM    = add.Command("ucspm", "Add a UCSPM Domain")
	updateUCSPM = update.Command("ucspm", "Update a UCSPM Domain")
	deleteUCSPM = delete.Command("ucspm", "Delete a UCSPM Domain")
	showUCSPM   = show.Command("ucspm", "Show a UCS Performance Manager")

	addUCSIP       = addUCS.Flag("ip", "IP Address or DNS name for UCS Manager, without http(s).").Required().IP()
	addUCSUsername = addUCS.Flag("username", "Name of user.").Required().String()
	addUCSPassword = addUCS.Flag("password", "Password for user in plain text.").Required().String()

	updateUCSIP       = updateUCS.Flag("ip", "IP Address or DNS name for UCS Manager, without http(s).").Required().IP()
	updateUCSUsername = updateUCS.Flag("username", "Name of user.").Required().String()
	updateUCSPassword = updateUCS.Flag("password", "Password for user in plain text.").Required().String()

	deleteUCSIP = deleteUCS.Flag("ip", "IP Address or DNS name for UCS Manager, without http(s).").Required().IP()

	showUCSIP = showUCS.Flag("ip", "IP Address or DNS name for UCS Manager, without http(s).").Required().IP()

	addUCSPMIP       = addUCSPM.Flag("ip", "IP Address or DNS name for UCS Performance Manager, without http(s).").Required().IP()
	addUCSPMUsername = addUCSPM.Flag("username", "Name of user.").Required().String()
	addUCSPMPassword = addUCSPM.Flag("password", "Password for user in plain text.").Required().String()

	updateUCSPMIP       = updateUCSPM.Flag("ip", "IP Address or DNS name for UCS Performance Manager, without http(s).").Required().IP()
	updateUCSPMUsername = updateUCSPM.Flag("username", "Name of user.").Required().String()
	updateUCSPMPassword = updateUCSPM.Flag("password", "Password for user in plain text.").Required().String()

	outputFile = output.Flag("set", "Configure the output filename, where the UUID and serial numbers will be saved.").Required().String()
	inputFile  = input.Flag("set", "Configure the input filename, where the UUID will be read from.").Required().String()

	runMonth = run.Flag("month", "Month process utilisation for").String()
	runYear  = run.Flag("year", "Year to process utilisation for").String()
)

func ProcessCommandLineArguments() string {
	switch kingpin.Parse() {
	case "run":
		return "RUN|" + *runMonth + "|" + *runYear
	case "add ucs":
		return "ADDUCS|" + as.ToString(*addUCSIP) + "|" + *addUCSUsername + "|" + *addUCSPassword
	case "update ucs":
		return "UPDATEUCS|" + as.ToString(*updateUCSIP) + "|" + *updateUCSUsername + "|" + *updateUCSPassword
	case "delete ucs":
		return "DELETEUCS|" + as.ToString(*deleteUCSIP)
	case "show ucs":
		return "SHOWUCS|" + as.ToString(*showUCSIP)
	case "show all":
		return "SHOWALL"
	case "add ucspm":
		return "ADDUCSPM|" + as.ToString(*addUCSPMIP) + "|" + *addUCSPMUsername + "|" + *addUCSPMPassword
	case "update ucspm":
		return "UPDATEUCSPM|" + as.ToString(*updateUCSPMIP) + "|" + *updateUCSPMUsername + "|" + *updateUCSPMPassword
	case "delete ucspm":
		return "DELETEUCSPM"
	case "show ucspm":
		return "SHOWUCSPM"
	case "input":
		return "SETINPUT|" + *inputFile
	case "output":
		return "SETOUTPUT|" + *outputFile
	}
	return ""
}
