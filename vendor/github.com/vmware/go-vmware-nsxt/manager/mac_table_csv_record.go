/* Copyright © 2017 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause

   Generated by: https://github.com/swagger-api/swagger-codegen.git */

package manager

type MacTableCsvRecord struct {

	// The MAC address
	MacAddress string `json:"mac_address"`

	// The virtual tunnel endpoint IP address
	VtepIp string `json:"vtep_ip,omitempty"`

	// The virtual tunnel endpoint MAC address
	VtepMacAddress string `json:"vtep_mac_address,omitempty"`
}
