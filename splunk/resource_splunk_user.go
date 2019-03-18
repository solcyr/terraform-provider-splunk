package splunk

import (
    "log"
    "github.com/hashicorp/terraform/helper/schema"
    "github.com/oliveagle/jsonpath"
    "fmt"
    "encoding/json"
    "net/url"
)

func resourceSplunkUser() *schema.Resource {
        return &schema.Resource{
                Create: resourceSplunkUserCreate,
                Read:   resourceSplunkUserRead,
                Update: resourceSplunkUserUpdate,
                Delete: resourceSplunkUserDelete,
                Importer: &schema.ResourceImporter{
                        State: schema.ImportStatePassthrough,
                },

                Schema: map[string]*schema.Schema{
                        "name": {
                                ForceNew: true,
                                Type:     schema.TypeString,
                                Required: true,
                        },
                        "password": {
                                Type:     schema.TypeString,
                                Required: true,
                        },
                        "real_name": {
                                Type:     schema.TypeString,
                                Optional: true,
                        },
                        "email": {
                                Type:     schema.TypeString,
                                Optional: true,
                        },
                        "roles": {
				Type:     schema.TypeList,
                                Optional: true,
                                Elem: &schema.Schema{
                                  Type:         schema.TypeString,
                                },
                        },
                },
        }
}

func resourceSplunkUserCreate(d *schema.ResourceData, meta interface{}) error {
        c := meta.(*Client)

	r := url.Values{}
        r.Set("name",     d.Get("name").(string))
        r.Set("password", d.Get("password").(string))
        r.Set("email",    d.Get("email").(string))
        r.Set("realname", d.Get("real_name").(string))
        r.Set("force-change-pass", "true")

        for _, element := range d.Get("roles").([]interface{}) {
            r.Add("roles", element.(string))
        }

        d.SetId(d.Get("name").(string))

        log.Printf("[DEBUG] Splunk User Creation: %s", d.Id())
        _, err := c.Post(PathUserCreate, r)
        if  err != nil  {
            return err
        }

        return resourceSplunkUserRead(d, meta)
}

func resourceSplunkUserRead(d *schema.ResourceData, meta interface{}) error {
        c := meta.(*Client)

        responseTxt, err := c.Get(fmt.Sprintf(PathUserSearch, url.QueryEscape(d.Id())))
        if err != nil {
            return err
        }

        var data interface{}
        json.Unmarshal([]byte(responseTxt), &data)

        res, err := jsonpath.JsonPathLookup(data, "$.entry[0].name")
        if err != nil {
            return err
        }
        d.SetId(res.(string))
        d.Set("name", res.(string))

        res, err = jsonpath.JsonPathLookup(data, "$.entry[0].content.email")
        if err != nil {
            return err
        }
        d.Set("email", res.(string))

        res, err = jsonpath.JsonPathLookup(data, "$.entry[0].content.realname")
        if err != nil {
            return err
        }
        d.Set("real_name", res.(string))

        res, err = jsonpath.JsonPathLookup(data, "$.entry[0].content.roles")
        if err != nil {
            return err
        }
        t := res.([]interface{})
        s := make([]string, len(t))
        for i, v := range t {
            s[i] = fmt.Sprint(v)
        }
        d.Set("roles", s)

        log.Printf("[DEBUG] Splunk User Read: %s", d.Get("name").(string))

        return err
}

func resourceSplunkUserUpdate(d *schema.ResourceData, meta interface{}) error {
        c := meta.(*Client)

	r := url.Values{}
        //r.Set("name",     d.Get("name").(string))
        r.Set("password", d.Get("password").(string))
        r.Set("email",    d.Get("email").(string))
        r.Set("realname", d.Get("real_name").(string))
        if d.HasChange("password") {
            r.Set("force-change-pass", "true")
        }

        for _, element := range d.Get("roles").([]interface{}) {
            r.Add("roles", element.(string))
        }

        log.Printf("[DEBUG] Splunk User Update: %s", d.Get("name").(string))
        _, err := c.Post(fmt.Sprintf(PathUserSearch, url.QueryEscape(d.Id())), r)
        if  err != nil  {
            return err
        }

        return resourceSplunkUserRead(d, meta)
}

func resourceSplunkUserDelete(d *schema.ResourceData, meta interface{}) error {
        c := meta.(*Client)

        log.Printf("[DEBUG] Splunk User Deletion: %s", d.Id())
        err := c.Delete(fmt.Sprintf(PathUserSearch, url.QueryEscape(d.Id())))

        return err

}

