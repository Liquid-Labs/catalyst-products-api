package products

import (
  "regexp"

  "github.com/Liquid-Labs/catalyst-core-api/go/resources/entities"
  "github.com/Liquid-Labs/go-nullable-mysql/nulls"
)

var phoneOutFormatter *regexp.Regexp = regexp.MustCompile(`^(\d{3})(\d{3})(\d{4})$`)

// On summary, we don't include address. Note leaving it empty and using
// 'omitempty' on the Product struct won't work because then Products without an address
// will appear 'incomplete' in the front-end model and never resolve.
type Product struct {
  entities.Entity
  LegalOwnerPubID nulls.String `json:"legalOwnerPubID"`
  DisplayName     nulls.String `json:"displayName"`
  Summary         nulls.String `json:"summary"`
  SupportEmail    nulls.String `json:"supportEmail"`
  SupportPhone    nulls.String `json:"supportPhone,string"`
  Homepage        nulls.String `json:"homepage"`
  LogoURL         nulls.String `json:"logoURL"`
  RepoURL         nulls.String `json:"repoURL"`
  IssuesURL       nulls.String `json:"issuesURL"`
  Ontology        nulls.String `json:"ontology"`
}

func (p *Product) FormatOut() {
  p.SupportPhone.String = phoneOutFormatter.ReplaceAllString(p.SupportPhone.String, `$1-$2-$3`)
}

func (p *Product) SetLegalOwnerPubID(val string) {
  p.LegalOwnerPubID = nulls.NewString(val)
}

func (p *Product) SetDisplayName(val string) {
  p.DisplayName = nulls.NewString(val)
}

func (p *Product) SetSummary(val string) {
  p.Summary = nulls.NewString(val)
}

func (p *Product) SetSupportEmail(val string) {
  p.SupportEmail = nulls.NewString(val)
}

func (p *Product) SetSupportPhone(val string) {
  p.SupportPhone = nulls.NewString(val)
}

func (p *Product) SetHomepage(val string) {
  p.Homepage = nulls.NewString(val)
}

func (p *Product) SetLogoURL(val string) {
  p.LogoURL = nulls.NewString(val)
}

func (p *Product) SetRepoURL(val string) {
  p.RepoURL = nulls.NewString(val)
}

func (p *Product) SetIssuesURL(val string) {
  p.IssuesURL = nulls.NewString(val)
}

func (p *Product) SetOntology(val string) {
  p.Ontology = nulls.NewString(val)
}

func (p *Product) Clone() *Product {
  return &Product{
    *p.Entity.Clone(),
    p.LegalOwnerPubID,
    p.DisplayName,
    p.Summary,
    p.SupportEmail,
    p.SupportPhone,
    p.Homepage,
    p.LogoURL,
    p.RepoURL,
    p.IssuesURL,
    p.Ontology,
  }
}
