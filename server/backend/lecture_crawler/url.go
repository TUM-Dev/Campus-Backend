package lecture_crawler

import "fmt"

func BuildLecturesListURL(pageNumber int) string {
	return fmt.Sprintf("%s/WBMODHB.cbShowMHBListe/NC_6180?pCaller=tabIdOrgModules&pOrgNr=1&pPageNr=%d", BaseURL, pageNumber)
}
