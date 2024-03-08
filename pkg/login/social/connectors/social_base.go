package connectors

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/oauth2"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/grafana/grafana/pkg/build/stringutil"
	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/grafana/grafana/pkg/login/social"
	"github.com/grafana/grafana/pkg/models/roletype"
	"github.com/grafana/grafana/pkg/services/auth/identity"
	"github.com/grafana/grafana/pkg/services/featuremgmt"
	"github.com/grafana/grafana/pkg/services/org"
	"github.com/grafana/grafana/pkg/services/ssosettings/validation"
	"github.com/grafana/grafana/pkg/setting"
	"github.com/grafana/grafana/pkg/util"
)

type SocialBase struct {
	*oauth2.Config
	orgService  org.Service
	info        *social.OAuthInfo
	cfg         *setting.Cfg
	reloadMutex sync.RWMutex
	log         log.Logger
	features    featuremgmt.FeatureToggles
}

func newSocialBase(name string,
	orgService org.Service,
	info *social.OAuthInfo,
	features featuremgmt.FeatureToggles,
	cfg *setting.Cfg,
) *SocialBase {
	logger := log.New("oauth." + name)

	return &SocialBase{
		Config:     createOAuthConfig(info, cfg, name),
		orgService: orgService,
		info:       info,
		log:        logger,
		features:   features,
		cfg:        cfg,
	}
}

func (s *SocialBase) updateInfo(name string, info *social.OAuthInfo) {
	s.Config = createOAuthConfig(info, s.cfg, name)
	s.info = info
}

type groupStruct struct {
	Groups []string `json:"groups"`
}

func (s *SocialBase) SupportBundleContent(bf *bytes.Buffer) error {
	s.reloadMutex.RLock()
	defer s.reloadMutex.RUnlock()

	return s.getBaseSupportBundleContent(bf)
}

func (s *SocialBase) GetOAuthInfo() *social.OAuthInfo {
	s.reloadMutex.RLock()
	defer s.reloadMutex.RUnlock()

	return s.info
}

func (s *SocialBase) AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string {
	s.reloadMutex.RLock()
	defer s.reloadMutex.RUnlock()

	return s.Config.AuthCodeURL(state, opts...)
}

func (s *SocialBase) Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	s.reloadMutex.RLock()
	defer s.reloadMutex.RUnlock()

	return s.Config.Exchange(ctx, code, opts...)
}

func (s *SocialBase) Client(ctx context.Context, t *oauth2.Token) *http.Client {
	s.reloadMutex.RLock()
	defer s.reloadMutex.RUnlock()

	return s.Config.Client(ctx, t)
}

func (s *SocialBase) TokenSource(ctx context.Context, t *oauth2.Token) oauth2.TokenSource {
	s.reloadMutex.RLock()
	defer s.reloadMutex.RUnlock()

	return s.Config.TokenSource(ctx, t)
}

func (s *SocialBase) getBaseSupportBundleContent(bf *bytes.Buffer) error {
	bf.WriteString("## Client configuration\n\n")
	bf.WriteString("```ini\n")
	bf.WriteString(fmt.Sprintf("allow_assign_grafana_admin = %v\n", s.info.AllowAssignGrafanaAdmin))
	bf.WriteString(fmt.Sprintf("allow_sign_up = %v\n", s.info.AllowSignup))
	bf.WriteString(fmt.Sprintf("allowed_domains = %v\n", s.info.AllowedDomains))
	bf.WriteString(fmt.Sprintf("auto_assign_org_role = %v\n", s.cfg.AutoAssignOrgRole))
	bf.WriteString(fmt.Sprintf("role_attribute_path = %v\n", s.info.RoleAttributePath))
	bf.WriteString(fmt.Sprintf("role_attribute_strict = %v\n", s.info.RoleAttributeStrict))
	bf.WriteString(fmt.Sprintf("skip_org_role_sync = %v\n", s.info.SkipOrgRoleSync))
	bf.WriteString(fmt.Sprintf("client_id = %v\n", s.Config.ClientID))
	bf.WriteString(fmt.Sprintf("client_secret = %v ; issue if empty\n", strings.Repeat("*", len(s.Config.ClientSecret))))
	bf.WriteString(fmt.Sprintf("auth_url = %v\n", s.Config.Endpoint.AuthURL))
	bf.WriteString(fmt.Sprintf("token_url = %v\n", s.Config.Endpoint.TokenURL))
	bf.WriteString(fmt.Sprintf("auth_style = %v\n", s.Config.Endpoint.AuthStyle))
	bf.WriteString(fmt.Sprintf("redirect_url = %v\n", s.Config.RedirectURL))
	bf.WriteString(fmt.Sprintf("scopes = %v\n", s.Config.Scopes))
	bf.WriteString("```\n\n")

	return nil
}

func (s *SocialBase) extractRoleAndAdminOptional(rawJSON []byte, groups []string) (org.RoleType, bool, error) {
	if s.info.RoleAttributePath == "" {
		if s.info.RoleAttributeStrict {
			return "", false, errRoleAttributePathNotSet.Errorf("role_attribute_path not set and role_attribute_strict is set")
		}
		return "", false, nil
	}

	if role, gAdmin := s.searchRole(rawJSON, groups); role.IsValid() {
		return role, gAdmin, nil
	} else if role != "" {
		return "", false, errInvalidRole.Errorf("invalid role: %s", role)
	}

	if s.info.RoleAttributeStrict {
		return "", false, errRoleAttributeStrictViolation.Errorf("idP did not return a role attribute, but role_attribute_strict is set")
	}

	return "", false, nil
}

func (s *SocialBase) extractRoleAndAdmin(rawJSON []byte, groups []string) (org.RoleType, bool, error) {
	role, gAdmin, err := s.extractRoleAndAdminOptional(rawJSON, groups)
	if role == "" {
		role = s.defaultRole()
	}

	return role, gAdmin, err
}

func (s *SocialBase) searchRole(rawJSON []byte, groups []string) (org.RoleType, bool) {
	role, err := util.SearchJSONForStringAttr(s.info.RoleAttributePath, rawJSON)
	if err == nil && role != "" {
		return getRoleFromSearch(role)
	}

	if groupBytes, err := json.Marshal(groupStruct{groups}); err == nil {
		role, err := util.SearchJSONForStringAttr(s.info.RoleAttributePath, groupBytes)
		if err == nil && role != "" {
			return getRoleFromSearch(role)
		}
	}

	return "", false
}

func (s *SocialBase) extractOrgRoles(ctx context.Context, rawJSON []byte, groups []string) (map[int64]org.RoleType, error) {
	if s.info.OrgMapping != nil && s.info.OrgAttributePath != "" {
		orgRoles, err := s.extractOrgRolesFromRaw(ctx, rawJSON, true)
		if err != nil || len(orgRoles) > 0 {
			return orgRoles, err
		}
		if groupBytes, err := json.Marshal(groupStruct{groups}); err == nil {
			return s.extractOrgRolesFromRaw(ctx, groupBytes, false)
		}
	}

	return nil, nil
}

func (s *SocialBase) extractOrgRolesFromRaw(ctx context.Context, rawJSON []byte, skipEmptyGroups bool) (map[int64]org.RoleType, error) {
	groups, err := util.SearchJSONForStringSliceAttr(s.info.OrgAttributePath, rawJSON)
	if err != nil {
		return nil, err
	}
	if skipEmptyGroups && len(groups) == 0 {
		return nil, nil
	}

	orgRoleMappings, err := s.resolveOrgMapping(groups)
	if err != nil {
		return nil, err
	}

	orgRoles := make(map[int64]org.RoleType, 0)
	for _, orgRoleMapping := range orgRoleMappings {

		if orgRoleMapping.OrgID == 0 && orgRoleMapping.OrgName != "" {
			getOrgQuery := &org.GetOrgByNameQuery{Name: orgRoleMapping.OrgName}
			res, err := s.orgService.GetByName(ctx, getOrgQuery)

			if err != nil {
				// ignore not existing org
				s.log.Warn("Unknown organization. Skipping.", "config_option", s.info.OrgAttributePath, "mapping", fmt.Sprintf("%v", orgRoleMapping))
				continue
			}
			orgRoleMapping.OrgID = res.ID
		} else if orgRoleMapping.OrgID <= 0 {
			s.log.Warn("Incorrect mapping found. Skipping.", "config_option", s.info.OrgAttributePath, "mapping", fmt.Sprintf("%v", orgRoleMapping))
			continue
		}
		orgRoles[orgRoleMapping.OrgID] = orgRoleMapping.Role
	}

	return orgRoles, nil
}

type OrgRoleMapping struct {
	OrgID   int64             `json:",string"`
	OrgName string            `type:"string"`
	Role    roletype.RoleType `type:"string" required:"true"`
}

func (s *SocialBase) resolveOrgMapping(groups []string) ([]OrgRoleMapping, error) {
	orgRoleMappings := make([]OrgRoleMapping, 0, len(s.info.OrgMapping))

	for _, m := range s.info.OrgMapping {
		parts := strings.SplitN(m, ":", 3)
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid tag format, expected 3 parts but got %d", len(parts))
		}

		if stringutil.Contains(groups, parts[0]) || parts[0] == "*" {
			if !roletype.RoleType(parts[2]).IsValid() {
				return nil, fmt.Errorf("invalid role type: %s", parts[2])
			}

			if orgId, err := strconv.Atoi(parts[1]); err == nil {
				orgRoleMappings = append(orgRoleMappings, OrgRoleMapping{
					OrgID: int64(orgId),
					Role:  roletype.RoleType(parts[2]),
				})
			} else {
				orgRoleMappings = append(orgRoleMappings, OrgRoleMapping{
					OrgName: parts[1],
					Role:    roletype.RoleType(parts[2]),
				})
			}
		}
	}

	return orgRoleMappings, nil
}

// defaultRole returns the default role for the user based on the autoAssignOrgRole setting
// if legacy is enabled "" is returned indicating the previous role assignment is used.
func (s *SocialBase) defaultRole() org.RoleType {
	if s.cfg.AutoAssignOrgRole != "" {
		s.log.Debug("No role found, returning default.")
		return org.RoleType(s.cfg.AutoAssignOrgRole)
	}

	// should never happen
	return org.RoleViewer
}

func (s *SocialBase) isGroupMember(groups []string) bool {
	if len(s.info.AllowedGroups) == 0 {
		return true
	}

	for _, allowedGroup := range s.info.AllowedGroups {
		for _, group := range groups {
			if group == allowedGroup {
				return true
			}
		}
	}

	return false
}

func (s *SocialBase) retrieveRawIDToken(idToken any) ([]byte, error) {
	tokenString, ok := idToken.(string)
	if !ok {
		return nil, fmt.Errorf("id_token is not a string: %v", idToken)
	}

	jwtRegexp := regexp.MustCompile("^([-_a-zA-Z0-9=]+)[.]([-_a-zA-Z0-9=]+)[.]([-_a-zA-Z0-9=]+)$")
	matched := jwtRegexp.FindStringSubmatch(tokenString)
	if matched == nil {
		return nil, fmt.Errorf("id_token is not in JWT format: %s", tokenString)
	}

	rawJSON, err := base64.RawURLEncoding.DecodeString(matched[2])
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding id_token: %w", err)
	}

	headerBytes, err := base64.RawURLEncoding.DecodeString(matched[1])
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding header: %w", err)
	}

	var header map[string]any
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return nil, fmt.Errorf("error deserializing header: %w", err)
	}

	if compressionVal, exists := header["zip"]; exists {
		compression, ok := compressionVal.(string)
		if !ok {
			return nil, fmt.Errorf("unrecognized compression header: %v", compressionVal)
		}

		if compression != "DEF" {
			return nil, fmt.Errorf("unknown compression algorithm: %s", compression)
		}

		fr, err := zlib.NewReader(bytes.NewReader(rawJSON))
		if err != nil {
			return nil, fmt.Errorf("error creating zlib reader: %w", err)
		}
		defer func() {
			if err := fr.Close(); err != nil {
				s.log.Warn("Failed closing zlib reader", "error", err)
			}
		}()

		rawJSON, err = io.ReadAll(fr)
		if err != nil {
			return nil, fmt.Errorf("error decompressing payload: %w", err)
		}
	}

	return rawJSON, nil
}

// match grafana admin role and translate to org role and bool.
// treat the JSON search result to ensure correct casing.
func getRoleFromSearch(role string) (org.RoleType, bool) {
	if strings.EqualFold(role, social.RoleGrafanaAdmin) {
		return org.RoleAdmin, true
	}

	return org.RoleType(cases.Title(language.Und).String(role)), false
}

func validateInfo(info *social.OAuthInfo, requester identity.Requester) error {
	return validation.Validate(info, requester,
		validation.RequiredValidator(info.ClientId, "Client Id"),
		validation.AllowAssignGrafanaAdminValidator,
		validation.SkipOrgRoleSyncAllowAssignGrafanaAdminValidator)
}
