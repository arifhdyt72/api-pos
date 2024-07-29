package external

import "github.com/gin-gonic/gin"

func CreateReportPDF(c *gin.Context) {
	// pdf := gofpdf.New("P", "mm", "A4", "")
	// pdf.AddPage()
	// pdf.SetFont("Arial", "B", 16)

	// pdf.Cell(190, 10, "List Transaction")
	// date := time.Now().Format("Mon, 02 January 2006")
	// pdf.Ln(10)
	// pdf.Cell(40, 10, date)
	// pdf.Ln(12)

	// // Headers
	// pdf.SetFont("Arial", "B", 12)
	// pdf.Cell(30, 10, "ID: "+helper.ToString(trx.ID))
	// pdf.Ln(7)
	// pdf.Cell(60, 10, "Transaction Date: "+trx.CreatedAt.Format("2006-01-02 15:04:05"))
	// pdf.Ln(7)
	// pdf.Cell(60, 10, "Amount: "+helper.ToString(trx.Total))
	// pdf.Ln(10)

	// pdf.SetFont("Arial", "B", 10)
	// pdf.Cell(80, 10, "Name")
	// pdf.Cell(40, 10, "Price")
	// pdf.Cell(40, 10, "Qty")
	// pdf.Cell(40, 10, "Subtotal")
	// pdf.Ln(10)
	// for _, val := range trx.TransactionDetails {
	// 	pdf.SetFont("Arial", "I", 10)
	// 	pdf.Cell(80, 10, helper.ToString(val.Item.Name))
	// 	pdf.Cell(40, 10, helper.ToString(val.Item.Price))
	// 	pdf.Cell(40, 10, helper.ToString(val.Qty))
	// 	pdf.Cell(40, 10, helper.ToString(val.Subtotal))
	// 	pdf.Ln(10)
	// }

	// yAfterCell := pdf.GetY() + 7
	// // Draw a line just below the cell
	// pdf.Line(10, yAfterCell, 200, yAfterCell)

	// err := pdf.OutputFileAndClose("alignedText.pdf")
	// if err != nil {
	// 	return err
	// }
}
