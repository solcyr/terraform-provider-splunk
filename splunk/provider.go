package splunk

import (
    "os"
    "github.com/hashicorp/terraform/helper/schema"
    "github.com/hashicorp/terraform/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
    return &schema.Provider{
        Schema: map[string]*schema.Schema{
            "url": &schema.Schema{
                Type:        schema.TypeString,
                Required:    true,
                DefaultFunc: schema.EnvDefaultFunc("SPLUNK_URL", nil),
                Description: "URL endpoint for Splunk API",
            },

            "username": &schema.Schema{
                Type:        schema.TypeString,
                Required:    true,
                DefaultFunc: schema.EnvDefaultFunc("SPLUNK_USERNAME", nil),
                Description: "The username for Splunk API operations.",
            },

            "password": &schema.Schema{
                Type:        schema.TypeString,
                Required:    true,
                DefaultFunc: schema.EnvDefaultFunc("SPLUNK_PASSWORD", nil),
                Description: "The password for Splunk API operations.",
            },

            "insecure_skip_verify": &schema.Schema{
                Type:        schema.TypeBool,
                Optional:    true,
                DefaultFunc: schema.EnvDefaultFunc("SPLUNK_INSECURE", false),
                Description: "Ignore certificate on Splunk server.",
            },
        },

        ResourcesMap: map[string]*schema.Resource{
            "splunk_saved_search": resourceSplunkSavedSearch(),
            "splunk_user": resourceSplunkUser(),
            "splunk_role": resourceSplunkRole(),
        },

        ConfigureFunc: providerConfigure,
    }
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

    username := d.Get("username").(string)
    if (password == "")  {
        username = os.Getenv("SPLUNK_USERNAME")
    }

    password := d.Get("password").(string)
    if (password == "")  {
        password = os.Getenv("SPLUNK_PASSWORD")
    }

    url := d.Get("url").(string)
    if (password == "")  {
        url = os.Getenv("SPLUNK_URL")
    }

    insecure_skip_verify := d.Get("insecure_skip_verify").(bool)
    if (password == "")  {
        insecure_skip_verify = os.Getenv("SPLUNK_INSECURE_SKIP_VERIFY")
    }

    config := Config{
        URL:                url,
        Username:           username,
        Password:           password,
        InsecureSkipVerify: insecure_skip_verify,
    }

    return config.Client()
}
