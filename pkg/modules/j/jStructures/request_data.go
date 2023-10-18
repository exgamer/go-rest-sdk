package jStructures

type RequestData struct {
	UserId           int
	UserType         string
	RequestId        string
	LanguageCode     string
	CompanyId        int
	CompanyIds       []int
	CurrentCompanyId int
	RequestMethod    string
	RequestUrl       string
}
