---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "grafana_folders Data Source - terraform-provider-grafana"
subcategory: ""
description: |-
  Datasource for retrieving all Grafana folders.
  Official documentation https://grafana.com/docs/grafana/latest/dashboards/dashboard_folders/Folder/Dashboard Search HTTP API https://grafana.com/docs/grafana/latest/http_api/folder_dashboard_search/Dashboard HTTP API https://grafana.com/docs/grafana/latest/http_api/dashboard/
---

# grafana_folders (Data Source)

Datasource for retrieving all Grafana folders.

* [Official documentation](https://grafana.com/docs/grafana/latest/dashboards/dashboard_folders/)
* [Folder/Dashboard Search HTTP API](https://grafana.com/docs/grafana/latest/http_api/folder_dashboard_search/)
* [Dashboard HTTP API](https://grafana.com/docs/grafana/latest/http_api/dashboard/)

## Example Usage

```terraform
resource "grafana_folder" "data_source_folders1" {
  title = "data_source_folders1"
}

resource "grafana_folder" "data_source_folders2" {
  title = "data_source_folders2"
}

// wait for folder resources to be created before searching
data "grafana_folders" "all" {
  depends_on = [
    grafana_folder.data_source_folders1,
    grafana_folder.data_source_folders2,
  ]
}

data "grafana_folders" "one" {
  limit = 1
  depends_on = [
    grafana_folder.data_source_folders1,
    grafana_folder.data_source_folders2,
  ]
}

// test to make sure search worked
data "grafana_folder" "test" {
  uid = data.grafana_folders.all.folders["data_source_folders1"]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- **id** (String) The ID of this resource.
- **limit** (Number) Maximum number of folders to return. Defaults to `1000`.

### Read-Only

- **folders** (Map of String)


