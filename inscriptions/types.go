package inscriptions

type InscriptionResult struct {
	InscriptionId string
	Tx            string
	Fees          int64
	ServiceFees   int64
}

type InscriptionResultRaw struct {
	Commit       string `json:"commit"`
	Inscriptions []struct {
		Id string `json:"id"`
	} `json:"inscriptions"`
}
