package catalog

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/oliy/pchome-cli/pkg/pchome/prodapi"
)

var productRefRe = regexp.MustCompile(`^[A-Z0-9]+-[A-Z0-9]+(?:-\d{3})?$`)

func NormalizeProductRef(ref string) (string, error) {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return "", fmt.Errorf("product id or PChome product URL is required")
	}

	if strings.Contains(ref, "://") {
		u, err := url.Parse(ref)
		if err != nil {
			return "", fmt.Errorf("invalid product reference %q", ref)
		}
		parts := strings.Split(strings.Trim(u.Path, "/"), "/")
		if len(parts) == 0 || parts[len(parts)-1] == "" {
			return "", fmt.Errorf("unable to extract product id from %q", ref)
		}
		ref = parts[len(parts)-1]
	}

	ref = strings.TrimSpace(strings.ToUpper(ref))
	ref = strings.TrimSuffix(ref, "/")
	ref = prodapi.NormalizeID(ref)

	if !productRefRe.MatchString(ref) {
		return "", fmt.Errorf("invalid product reference %q", ref)
	}

	return ref, nil
}
