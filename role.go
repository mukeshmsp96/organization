package organization

import (
	"errors"
	"fmt"

	ldap "gopkg.in/ldap.v2"
)

// AddRole to ldap server, this method will automatically update org's rbacx
func (org *Organization) AddRole(name, description string, ups, pps []string) (string, error) {

	id := generateNewID()
	dn := org.dn(id, role)
	aq := ldap.NewAddRequest(dn)

	aq.Attribute(`objectClass`, []string{`role`, `top`})

	aq.Attribute(`cn`, []string{name})
	aq.Attribute(`description`, []string{description})
	aq.Attribute(`upid`, ups)
	aq.Attribute(`ppid`, pps)

	err := org.Add(aq) // 先写数据库
	if err != nil {
		return ``, err
	}

	r, err := org.RoleByID(id)
	if err != nil {
		return ``, nil
	}
	org.insertNewPolicyByRole(r)

	return id, nil
}

// DelRole from ldap server, automatically update org's rbacx
func (org *Organization) DelRole(id string) error {

	// 判断有没有人引用这个Role
	mIDs, err := org.MemberIDsByRoleIDs([]string{id})
	if err != nil {
		return err
	}
	if len(mIDs) > 0 {
		return fmt.Errorf(`尚有人引用此角色 count: %d`, len(mIDs))
	}

	dn := org.dn(id, role)
	dq := ldap.NewDelRequest(dn, nil)

	err = org.Del(dq)
	if err != nil {
		return err
	}

	org.removePolicyByRoleID(id)

	return nil
}

// ModifyRole in ldap server, automatically update org's rbacx
func (org *Organization) ModifyRole(id, name, description string, ups, pps []string) error {
	if len(ups) > 0 || len(pps) > 0 { // permission id 有效性判断
		p, err := org.PermissionByIDs(append(ups, pps...))
		if err != nil {
			return err
		}
		if len(p.Data) != len(ups)+len(pps) {
			return errors.New(`permission ids is invalid`)
		}
	}

	r, err := org.RoleByID(id)
	if err != nil {
		return err
	}

	mq := ldap.NewModifyRequest(r[`dn`].(string))

	if len(name) > 0 {
		mq.Replace(`cn`, []string{name})
	}
	if len(description) > 0 {
		mq.Replace(`description`, []string{description})
	}
	if len(ups) > 0 {
		mq.Replace(`upid`, ups)
	}
	if len(pps) > 0 {
		mq.Replace(`ppid`, pps)
	}

	err = org.Modify(mq)
	if err != nil {
		return err
	}

	oTypes := append(ups, r[`upid`].([]string)...)
	nTypes := append(pps, r[`ppid`].([]string)...)

	org.refreshRBACIfNeeded(oTypes, nTypes)

	return nil
}

// AllRoles ...
func (org *Organization) AllRoles() ([]map[string]interface{}, error) {
	r, e := org.Roles(0, nil)
	if e != nil {
		return nil, e
	}
	return r.Data, nil
}

// Roles in ldap
func (org *Organization) Roles(pageSize uint32, cookie []byte) (*SearchResult, error) {
	return org.searchRole(``, pageSize, cookie)
}

// RoleByID ...
func (org *Organization) RoleByID(id string) (map[string]interface{}, error) {
	sr, e := org.RoleByIDs([]string{id})
	if e != nil {
		return nil, e
	}
	if len(sr.Data) != 1 {
		return nil, errors.New(`found many roles`)
	}
	return sr.Data[0], nil
}

// RoleByIDs in ldap
func (org *Organization) RoleByIDs(ids []string) (*SearchResult, error) {
	filter, err := sqConvertIDsToFilter(ids)
	if err != nil {
		return nil, err
	}
	return org.searchRole(filter, 0, nil)
}

// RoleIDsByMemberID ...
func (org *Organization) RoleIDsByMemberID(id string) ([]string, error) {

	if len(id) == 0 {
		return nil, errors.New(`id must not be empty`)
	}

	filter := fmt.Sprintf(`(id=%s)`, id)

	sq := ldap.NewSearchRequest(org.parentDN(member),
		ldap.ScopeSingleLevel,
		ldap.DerefAlways, 0, 0, false, filter, []string{`rbacRole`}, nil)
	sr, err := org.Search(sq)
	if err != nil {
		return nil, err
	}
	if len(sr.Entries) != 1 {
		return nil, fmt.Errorf(`[%s] member doesn't exit`, id)
	}
	return sr.Entries[0].GetAttributeValues(`rbacRole`), nil
}

// RoleIDsByPermissionID which role contain this permission
func (org *Organization) RoleIDsByPermissionID(id string) ([]string, error) {

	filter := fmt.Sprintf(`(|(upid=%s)(ppid=%s))`, id, id)
	dn := org.parentDN(role)

	sq := &searchRequest{dn, filter, []string{`id`}, nil, 0, nil}

	r, e := org.search(sq)
	if e != nil {
		return nil, e
	}

	var ids []string
	for _, v := range r.Data {
		ids = append(ids, v[`id`].(string))
	}

	return ids, nil
}
