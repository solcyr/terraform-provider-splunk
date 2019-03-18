package splunk

import (
    "log"
    "github.com/hashicorp/terraform/helper/schema"
    "github.com/oliveagle/jsonpath"
    "fmt"
    "encoding/json"
    "net/url"
)

func resourceSplunkRole() *schema.Resource {
        return &schema.Resource{
                Create: resourceSplunkRoleCreate,
                Read:   resourceSplunkRoleRead,
                Update: resourceSplunkRoleUpdate,
                Delete: resourceSplunkRoleDelete,
                Importer: &schema.ResourceImporter{
                        State: schema.ImportStatePassthrough,
                },

                Schema: map[string]*schema.Schema{
                        "name": {
                                ForceNew: true,
                                Type:     schema.TypeString,
                                Required: true,
                        },
                        "search_filter": {
                                Type:     schema.TypeString,
                                Optional: true,
                        },
                        "imported_roles": {
				Type:     schema.TypeList,
                                Optional: true,
                                Elem: &schema.Schema{
                                  Type:         schema.TypeString,
                                },
                        },
                        "indexes_allowed": {
				Type:     schema.TypeList,
                                Optional: true,
                                Elem: &schema.Schema{
                                  Type:         schema.TypeString,
                                },
                        },
                },
        }
}

func resourceSplunkRoleCreate(d *schema.ResourceData, meta interface{}) error {
        c := meta.(*Client)

	r := url.Values{}
        r.Set("name",       d.Get("name").(string))
        r.Set("srchFilter", d.Get("search_filter").(string))
        r.Set("defaultApp", "search")

        for _, element := range d.Get("indexes_allowed").([]interface{}) {
            r.Add("srchIndexesAllowed", element.(string))
        }

        for _, element := range d.Get("imported_roles").([]interface{}) {
            r.Add("imported_roles", element.(string))
        }

        d.SetId(d.Get("name").(string))

        log.Printf("[DEBUG] Splunk Role Creation: %s", d.Id())
        _, err := c.Post(PathRoleCreate, r)
        if  err != nil  {
            return err
        }

        return resourceSplunkRoleRead(d, meta)
}

func resourceSplunkRoleRead(d *schema.ResourceData, meta interface{}) error {
        c := meta.(*Client)

        responseTxt, err := c.Get(fmt.Sprintf(PathRoleSearch, url.QueryEscape(d.Id())))
        if err != nil {
            return err
        }

        var data interface{}
        json.Unmarshal([]byte(responseTxt), &data)

        res, err := jsonpath.JsonPathLookup(data, "$.entry.name[0]")
        if err != nil {
            return err
        }
        d.SetId(res.(string))
        d.Set("name", res.(string))

        res, err = jsonpath.JsonPathLookup(data, "$.entry.content.srchFilter[0]")
        if err != nil {
            return err
        }
        d.Set("search_filter", res.(string))

        res, err = jsonpath.JsonPathLookup(data, "$.entry.content.srchIndexesAllowed[0]")
        if err != nil {
            return err
        }
        t := res.([]interface{})
        s := make([]string, len(t))
        for i, v := range t {
            s[i] = fmt.Sprint(v)
        }
        d.Set("indexes_allowed", s)

        res, err = jsonpath.JsonPathLookup(data, "$.entry.content.imported_roles[0]")
        if err != nil {
            return err
        }
        t = res.([]interface{})
        s = make([]string, len(t))
        for i, v := range t {
            s[i] = fmt.Sprint(v)
        }
        d.Set("imported_roles", s)

        log.Printf("[DEBUG] Splunk Role Read: %s", d.Get("name").(string))

        return err
}

func resourceSplunkRoleUpdate(d *schema.ResourceData, meta interface{}) error {
        c := meta.(*Client)

	r := url.Values{}
        //r.Set("name",       d.Get("name").(string))
        r.Set("srchFilter", d.Get("search_filter").(string))
        r.Set("defaultApp", "search")

        for _, element := range d.Get("indexes_allowed").([]interface{}) {
            r.Add("srchIndexesAllowed", element.(string))
        }

        for _, element := range d.Get("imported_roles").([]interface{}) {
            r.Add("imported_roles", element.(string))
        }

        log.Printf("[DEBUG] Splunk Role Update: %s", d.Get("name").(string))
        _, err := c.Post(fmt.Sprintf(PathRoleSearch, url.QueryEscape(d.Id())), r)
        if  err != nil  {
            return err
        }

        return resourceSplunkRoleRead(d, meta)
}

func resourceSplunkRoleDelete(d *schema.ResourceData, meta interface{}) error {
        c := meta.(*Client)

        log.Printf("[DEBUG] Splunk Role Deletion: %s", d.Id())
        err := c.Delete(fmt.Sprintf(PathRoleSearch, url.QueryEscape(d.Id())))

        return err

}

