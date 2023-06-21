package models

type CTM struct {
	ID       string `json:"id" bson:"_id,omitempty"`
	Email    string `json:"email,omitempty" bson:"email,omitempty"`
	Address  string `json:"address,omitempty" bson:"address,omitempty"`
	Business string `json:"business,omitempty" bson:"business,omitempty"`
	Name     string `json:"name,omitempty" bson:"name,omitempty"`
	Locale   string `json:"locale,omitempty" bson:"locale,omitempty"`
	Customer string `json:"customer,omitempty" bson:"customer,omitempty"`
}

type ADD struct {
	ID          string `json:"id,omitempty" bson:"_id,omitempty"`
	Status      string `json:"status,omitempty" bson:"status,omitempty"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
	FirstLine   string `json:"first_line,omitempty" bson:"first_line,omitempty"`
	City        string `json:"city,omitempty" bson:"city,omitempty"`
	PostalCode  string `json:"postal_code,omitempty" bson:"postal_code,omitempty"`
	Region      string `json:"region,omitempty" bson:"region,omitempty"`
	CountryCode string `json:"country_code,omitempty" bson:"country_code,omitempty"`
	Customer    string `json:"customer,omitempty" bson:"customer,omitempty"`
	PadID       string `json:"pad_id,omitempty" bson:"pad_id,omitempty"`
}

type BIZ struct {
	ID            string `json:"id,omitempty" bson:"_id,omitempty"`
	Status        string `json:"status,omitempty" bson:"status,omitempty"`
	Name          string `json:"name,omitempty" bson:"name,omitempty"`
	CompanyNumber string `json:"company_number,omitempty" bson:"company_number,omitempty"`
	TaxIdentifier string `json:"tax_identifier,omitempty" bson:"tax_identifier,omitempty"`
	Customer      string `json:"customer,omitempty" bson:"customer,omitempty"`
	PadID         string `json:"pad_id,omitempty" bson:"pad_id,omitempty"`
}
