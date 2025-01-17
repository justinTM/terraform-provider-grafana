package grafana

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DatasourceFolders() *schema.Resource {
	return &schema.Resource{
		Description: `
Datasource for retrieving all Grafana folders.

* [Official documentation](https://grafana.com/docs/grafana/latest/dashboards/dashboard_folders/)
* [Folder/Dashboard Search HTTP API](https://grafana.com/docs/grafana/latest/http_api/folder_dashboard_search/)
* [Dashboard HTTP API](https://grafana.com/docs/grafana/latest/http_api/dashboard/)
`,
		ReadContext: dataSourceReadFolders,
		Schema: map[string]*schema.Schema{
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1000,
				Description: "Maximum number of folders to return.",
			},
			"folders": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     schema.TypeString,
			},
		},
	}
}

func dataSourceReadFolders(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client).gapi
	var diags diag.Diagnostics
	params := url.Values{
		"limit": {fmt.Sprint(d.Get("limit"))},
		"type":  {"dash-folder"},
	}

	d.SetId(hashDashboardSearchParameters(params))

	results, err := client.FolderDashboardSearch(params)
	if err != nil {
		return diag.FromErr(err)
	}

	folders := make(map[string]interface{}, len(results))
	for _, result := range results {
		folders[result.Title] = result.UID
	}

	if err := d.Set("folders", folders); err != nil {
		return diag.Errorf("error setting folders attribute: %s", err)
	}

	return diags
}
