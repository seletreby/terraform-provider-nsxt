/* Copyright © 2017 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	api "github.com/vmware/go-vmware-nsxt"
	"github.com/vmware/go-vmware-nsxt/manager"
	"net/http"
)

func resourceAlgTypeNsService() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlgTypeNsServiceCreate,
		Read:   resourceAlgTypeNsServiceRead,
		Update: resourceAlgTypeNsServiceUpdate,
		Delete: resourceAlgTypeNsServiceDelete,

		Schema: map[string]*schema.Schema{
			"revision":     getRevisionSchema(),
			"system_owned": getSystemOwnedSchema(),
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Description of this resource",
				Optional:    true,
			},
			"display_name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Defaults to ID if not set",
				Optional:    true,
			},
			"tags": getTagsSchema(),
			"default_service": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "The default NSServices are created in the system by default. These NSServices can't be modified/deleted",
				Computed:    true,
			},
			"destination_ports": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Range of destination ports. This is single value, not a set",
				Required:    true,
			},
			"source_ports": &schema.Schema{
				Type:        schema.TypeSet,
				Description: "Set of source ports or ranges",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
			},
			"alg": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Algorithm",
				Required:    true,
			},
		},
	}
}

func resourceAlgTypeNsServiceCreate(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	description := d.Get("description").(string)
	display_name := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	default_service := d.Get("default_service").(bool)
	alg := d.Get("alg").(string)
	source_ports := getStringListFromSchemaSet(d, "source_ports")
	destination_ports := make([]string, 0, 1)
	destination_ports = append(destination_ports, d.Get("destination_ports").(string))

	ns_service := manager.AlgTypeNsService{
		NsService: manager.NsService{
			Description:    description,
			DisplayName:    display_name,
			Tags:           tags,
			DefaultService: default_service,
		},
		NsserviceElement: manager.AlgTypeNsServiceEntry{
			ResourceType:     "ALGTypeNSService",
			Alg:              alg,
			DestinationPorts: destination_ports,
			SourcePorts:      source_ports,
		},
	}

	ns_service, resp, err := nsxClient.GroupingObjectsApi.CreateAlgTypeNSService(nsxClient.Context, ns_service)

	if err != nil {
		return fmt.Errorf("Error during NsService create: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Unexpected status returned during NsService create: %v", resp.StatusCode)
	}
	d.SetId(ns_service.Id)
	return resourceAlgTypeNsServiceRead(d, m)
}

func resourceAlgTypeNsServiceRead(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining ns service id")
	}

	ns_service, resp, err := nsxClient.GroupingObjectsApi.ReadAlgTypeNSService(nsxClient.Context, id)
	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("NsService %s not found", id)
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error during NsService read: %v", err)
	}

	nsservice_element := ns_service.NsserviceElement

	d.Set("revision", ns_service.Revision)
	d.Set("system_owned", ns_service.SystemOwned)
	d.Set("description", ns_service.Description)
	d.Set("display_name", ns_service.DisplayName)
	setTagsInSchema(d, ns_service.Tags)
	d.Set("default_service", ns_service.DefaultService)
	d.Set("alg", nsservice_element.Alg)
	d.Set("destination_ports", nsservice_element.DestinationPorts)
	d.Set("source_ports", nsservice_element.SourcePorts)

	return nil
}

func resourceAlgTypeNsServiceUpdate(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining ns service id")
	}

	description := d.Get("description").(string)
	display_name := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	default_service := d.Get("default_service").(bool)
	alg := d.Get("alg").(string)
	source_ports := getStringListFromSchemaSet(d, "source_ports")
	destination_ports := make([]string, 0, 1)
	destination_ports = append(destination_ports, d.Get("destination_ports").(string))
	revision := int64(d.Get("revision").(int))

	ns_service := manager.AlgTypeNsService{
		NsService: manager.NsService{
			Description:    description,
			DisplayName:    display_name,
			Tags:           tags,
			DefaultService: default_service,
			Revision:       revision,
		},
		NsserviceElement: manager.AlgTypeNsServiceEntry{
			ResourceType:     "ALGTypeNSService",
			Alg:              alg,
			DestinationPorts: destination_ports,
			SourcePorts:      source_ports,
		},
	}

	ns_service, resp, err := nsxClient.GroupingObjectsApi.UpdateAlgTypeNSService(nsxClient.Context, id, ns_service)
	if err != nil || resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("Error during NsService update: %v %v", err, resp)
	}

	return resourceAlgTypeNsServiceRead(d, m)
}

func resourceAlgTypeNsServiceDelete(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining ns service id")
	}

	localVarOptionals := make(map[string]interface{})
	resp, err := nsxClient.GroupingObjectsApi.DeleteNSService(nsxClient.Context, id, localVarOptionals)
	if err != nil {
		return fmt.Errorf("Error during NsService delete: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("NsService %s not found", id)
		d.SetId("")
	}
	return nil
}