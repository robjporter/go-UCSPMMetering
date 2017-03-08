package eula

import (
	"fmt"
	"log"
)

func DisplayEULA() bool {
	fmt.Println(`
	Software License Agreement

	1. This is an agreement between Licensor and Licensee, who is being licensed to use the named Software.

	2. Licensee acknowledges that this is only a limited nonexclusive license. Licensor is and remains the owner of all titles, rights, and interests in the Software.

	3. This License permits Licensee to install the Software on more than one computer system, as long as the Software will not be used on more than one computer system simultaneously. Licensee will not make copies of the Software or allow copies of the Software to be made by others, unless authorized by this License Agreement. Licensee may make copies of the Software for backup purposes only.

	4. This Software is subject to a limited warranty. Licensor warrants to Licensee that the physical medium on which this Software is distributed is free from defects in materials and workmanship under normal use, the Software will perform according to its printed documentation, and to the best of Licensor's knowledge Licensee's use of this Software according to the printed documentation is not an infringement of any third party's intellectual property rights. This limited warranty lasts for a period of ____ days after delivery. To the extent permitted by law, THE ABOVE-STATED LIMITED WARRANTY REPLACES ALL OTHER WARRANTIES, EXPRESS OR IMPLIED, AND LICENSOR DISCLAIMS ALL IMPLIED WARRANTIES INCLUDING ANY IMPLIED WARRANTY OF TITLE, MERCHANTABILITY, NONINFRINGEMENT, OR OF FITNESS FOR A PARTICULAR PURPOSE. No agent of Licensor is authorized to make any other warranties or to modify this limited warranty. Any action for breach of this limited warranty must be commenced within one year of the expiration of the warranty. Because some jurisdictions do not allow any limit on the length of an implied warranty, the above limitation may not apply to this Licensee. If the law does not allow disclaimer of implied warranties, then any implied warranty is limited to ____ days after delivery of the Software to Licensee. Licensee has specific legal rights pursuant to this warranty and, depending on Licensee's jurisdiction, may have additional rights.

	5. In case of a breach of the Limited Warranty, Licensee's exclusive remedy is as follows: Licensee will return all copies of the Software to Licensor, at Licensee's cost, along with proof of purchase. (Licensee can obtain a step-by-step explanation of this procedure, including a return authorization code, by contacting Licensor at [address and toll free telephone number].) At Licensor's option, Licensor will either send Licensee a replacement copy of the Software, at Licensor's expense, or issue a full refund.

	6. Notwithstanding the foregoing, LICENSOR IS NOT LIABLE TO LICENSEE FOR ANY DAMAGES, INCLUDING COMPENSATORY, SPECIAL, INCIDENTAL, EXEMPLARY, PUNITIVE, OR CONSEQUENTIAL DAMAGES, CONNECTED WITH OR RESULTING FROM THIS LICENSE AGREEMENT OR LICENSEE'S USE OF THIS SOFTWARE. Licensee's jurisdiction may not allow such a limitation of damages, so this limitation may not apply.

	7. Licensee agrees to defend and indemnify Licensor and hold Licensor harmless from all claims, losses, damages, complaints, or expenses connected with or resulting from Licensee's business operations.

	8. Licensor has the right to terminate this License Agreement and Licensee's right to use this Software upon any material breach by Licensee.

	9. Licensee agrees to return to Licensor or to destroy all copies of the Software upon termination of the License.

	10. This License Agreement is the entire and exclusive agreement between Licensor and Licensee regarding this Software. This License Agreement replaces and supersedes all prior negotiations, dealings, and agreements between Licensor and Licensee regarding this Software.

	11. This License Agreement is governed by the law of [State] applicable to [State] contracts.

	12. This License Agreement is valid without Licensor's signature. It becomes effective upon the earlier of Licensee's signature or Licensee's use of the Software.`)

	fmt.Println("\nPress read and confirm acceptance with y/Y/yes/YES")

	output := askForConfirmation()

	return output

}

// askForConfirmation uses Scanln to parse user input. A user must type in "yes" or "no" and
// then press enter. It has fuzzy matching, so "y", "Y", "yes", "YES", and "Yes" all count as
// confirmations. If the input is not recognized, it will ask again. The function does not return
// until it gets a valid response from the user. Typically, you should use fmt to print out a question
// before calling askForConfirmation. E.g. fmt.Println("WARNING: Are you sure? (yes/no)")
func askForConfirmation() bool {
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}
	okayResponses := []string{"y", "Y", "yes", "Yes", "YES"}
	nokayResponses := []string{"n", "N", "no", "No", "NO"}
	if containsString(okayResponses, response) {
		return true
	} else if containsString(nokayResponses, response) {
		return false
	} else {
		fmt.Println("Please type yes or no and then press enter:")
		return askForConfirmation()
	}
}

// posString returns the first index of element in slice.
// If slice does not contain element, returns -1.
func posString(slice []string, element string) int {
	for index, elem := range slice {
		if elem == element {
			return index
		}
	}
	return -1
}

// containsString returns true iff slice contains element
func containsString(slice []string, element string) bool {
	return !(posString(slice, element) == -1)
}
