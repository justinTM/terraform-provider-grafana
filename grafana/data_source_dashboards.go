package grafana

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DatasourceDashboards() *schema.Resource {
	return &schema.Resource{
		Description: `
Datasource for retrieving all dashboards. Specify list of folder IDs to search in for dashboards.

* [Official documentation](https://grafana.com/docs/grafana/latest/dashboards/)
* [Folder/Dashboard Search HTTP API](https://grafana.com/docs/grafana/latest/http_api/folder_dashboard_search/)
* [Dashboard HTTP API](https://grafana.com/docs/grafana/latest/http_api/dashboard/)
`,
		ReadContext: dataSourceReadDashboards,
		Schema: map[string]*schema.Schema{
			"folder_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "Numerical IDs of Grafana folders containing dashboards. Specify to filter for dashboards by folder, or leave blank to get all dashboards.",
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `List of string Grafana dashboard tags to search for, eg. ["prod"]. Used only as search input, i.e., attribute value will remain unchanged.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"dashboards": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Map of Grafana dashboard unique identifiers (list of string UIDs as values) to folder IDs (integers as keys).",
				Elem: &schema.Schema{
					Type:        schema.TypeList,
					Description: "List of string dashboard UIDs.",
					Elem:        &schema.Schema{Type: schema.TypeString},
				},
			},
		},
	}
}

func dataSourceReadDashboards(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client).gapi
	params := map[string]string{
		"limit": "5000",
		"type":  "dash-db",
	}

	// add tags and folder IDs from attributes to dashboard search parameters
	for thisParamKey, thisListAttributeName := range map[string]string{
		"folderIds": "folder_ids",
		"tag":       "tags",
	} {
		list := d.Get(thisListAttributeName).([]interface{})
		if len(list) > 0 {
			listJSON, err := json.Marshal(list)
			if err != nil {
				return diag.FromErr(err)
			}
			params[thisParamKey] = string(listJSON)
		}
	}

	results, err := client.FolderDashboardSearch(params)
	if err != nil {
		return diag.FromErr(err)
	}

	// make list of string dashboard UIDs (as values) mapped to each int folder ID (as keys)
	dashboards := make(map[int64][]string, len(results))
	dashboardIDs := make([]int64, len(results))
	dashboardUIDs := make([]string, len(results))
	for i, thisResult := range results {
		thisFolderID := int64(thisResult.FolderID)
		dashboards[thisFolderID] = append(dashboards[thisFolderID], thisResult.UID)
		dashboardIDs[i] = int64(thisResult.ID)
		dashboardUIDs[i] = thisResult.UID
	}

	var folders []int64
	for thisFolderID := range dashboards {
		folders = append(folders, thisFolderID)
	}

	d.Set("dashboards", dashboards)

	return nil
}
