package interview

import "fmt"

/**
 * Companies have parent-child relationships
 * When a customer takes a loan with a company, the loan is issued under the topmost parent company
 */
func resolveTopLevelIssuers(companyHierarchy map[string][]string, loanCompanies []string) ([]string, error) {
	// Build child → parent lookup
	childToParent := map[string]string{}
	for parent, subsidiaries := range companyHierarchy {
		for _, subsidiary := range subsidiaries {
			childToParent[subsidiary] = parent
		}
	}

	// Walk up hierarchy to topmost ancestor
	findTopLevelIssuer := func(company string) (string, error) {
		visited := map[string]bool{}
		for {
			if visited[company] {
				return "", fmt.Errorf("cycle detected at: %s", company)
			}
			visited[company] = true

			parent, hasParent := childToParent[company]
			if !hasParent {
				return company, nil
			}
			company = parent
		}
	}

	var issuers []string
	for _, loanCompany := range loanCompanies {
		issuer, err := findTopLevelIssuer(loanCompany)
		if err != nil {
			return nil, err
		}
		issuers = append(issuers, issuer)
	}

	return issuers, nil
}
