// Copyright 2026 The MathWorks, Inc.

//go:build windows

package sessiondetails

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

func ensureFileSecure(path string) error {
	currentToken := windows.GetCurrentProcessToken()
	user, err := currentToken.GetTokenUser()
	if err != nil {
		return fmt.Errorf("failed to get current process token user: %w", err)
	}

	if user == nil || user.User.Sid == nil {
		return fmt.Errorf("current process token user SID is missing")
	}

	entries := []windows.EXPLICIT_ACCESS{
		{
			AccessPermissions: windows.GENERIC_ALL,
			AccessMode:        windows.SET_ACCESS,
			Inheritance:       windows.NO_INHERITANCE,
			Trustee: windows.TRUSTEE{
				TrusteeForm:  windows.TRUSTEE_IS_SID,
				TrusteeType:  windows.TRUSTEE_IS_USER,
				TrusteeValue: windows.TrusteeValueFromSID(user.User.Sid),
			},
		},
	}

	acl, err := windows.ACLFromEntries(entries, nil)
	if err != nil {
		return fmt.Errorf("failed to build ACL for session details file: %w", err)
	}

	if err := windows.SetNamedSecurityInfo(
		path,
		windows.SE_FILE_OBJECT,
		windows.DACL_SECURITY_INFORMATION|windows.PROTECTED_DACL_SECURITY_INFORMATION|windows.OWNER_SECURITY_INFORMATION,
		user.User.Sid,
		nil,
		acl,
		nil,
	); err != nil {
		return fmt.Errorf("failed to set security info on session details file: %w", err)
	}

	return nil
}

func AssertFileSecure(path string) error {
	securityDescriptor, err := windows.GetNamedSecurityInfo(
		path,
		windows.SE_FILE_OBJECT,
		windows.OWNER_SECURITY_INFORMATION|windows.DACL_SECURITY_INFORMATION,
	)
	if err != nil {
		return fmt.Errorf("failed to get owner security info: %w", err)
	}

	owner, _, err := securityDescriptor.Owner()
	if err != nil {
		return fmt.Errorf("failed to read owner SID: %w", err)
	}
	if owner == nil {
		return fmt.Errorf("owner SID is missing")
	}

	currentToken := windows.GetCurrentProcessToken()
	user, err := currentToken.GetTokenUser()
	if err != nil {
		return fmt.Errorf("failed to get current process token user: %w", err)
	}

	if user == nil || user.User.Sid == nil {
		return fmt.Errorf("current process token user SID is missing")
	}

	if !owner.Equals(user.User.Sid) {
		return fmt.Errorf("session details file owner does not match current user")
	}

	dacl, _, err := securityDescriptor.DACL()
	if err != nil {
		return fmt.Errorf("failed to read DACL: %w", err)
	}

	if dacl == nil {
		return fmt.Errorf("session details file DACL is missing")
	}

	if dacl.AceCount != 1 {
		return fmt.Errorf("expected exactly one ACE for current user, got %d", dacl.AceCount)
	}

	var ace *windows.ACCESS_ALLOWED_ACE
	if err := windows.GetAce(dacl, 0, &ace); err != nil {
		return fmt.Errorf("failed to read ACE from DACL: %w", err)
	}

	if ace.Header.AceType != windows.ACCESS_ALLOWED_ACE_TYPE {
		return fmt.Errorf("unexpected ACE type %d in DACL", ace.Header.AceType)
	}

	aceSID := (*windows.SID)(unsafe.Pointer(&ace.SidStart)) //nolint:gosec // Required to extract SID from ACCESS_ALLOWED_ACE
	if !aceSID.Equals(user.User.Sid) {
		return fmt.Errorf("DACL ACE does not target the current user")
	}

	return nil
}
