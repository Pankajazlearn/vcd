package vcd

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceVcdOrgUser() *schema.Resource {
	return &schema.Resource{
		Read: datasourceVcdOrgUserRead,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"name", "user_id"},
				Description:  `User's name. Required if "user_id" is not set`,
			},
			"user_id": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"name", "user_id"},
				Description:  `User's id. Required if "name" is not set`,
			},
			"org": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "The name of organization to use, optional if defined at provider " +
					"level. Useful when connected as sysadmin working across different organizations",
			},
			"role": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Role within the organization",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user's description",
			},
			"provider_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identity provider type for this this user. One of: 'INTEGRATED', 'SAML', 'OAUTH'. ",
			},
			"full_name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user's full name",
			},
			"email_address": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user's email address",
			},
			"telephone": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user's telephone",
			},
			"instant_messaging": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user's telephone",
			},
			"enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "True if the user is enabled and can log in.",
			},
			"is_group_role": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "True if this user has a group role.",
			},
			"is_locked": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
				Description: "If the user account has been locked due to too many invalid login attempts, the value " +
					"will change to true (only the system can lock the user). ",
			},
			"deployed_vm_quota": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Quota of vApps that this user can deploy. A value of 0 specifies an unlimited quota.",
			},
			"stored_vm_quota": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Quota of vApps that this user can store. A value of 0 specifies an unlimited quota.",
			},
		},
	}
}

func datasourceVcdOrgUserRead(d *schema.ResourceData, meta interface{}) error {

	var identifier string
	name := d.Get("name").(string)
	id := d.Get("user_id").(string)

	if name != "" {
		identifier = name
	} else {
		identifier = id
	}

	vcdClient := meta.(*VCDClient)
	adminOrg, err := vcdClient.GetAdminOrgFromResource(d)
	if err != nil {
		return fmt.Errorf("[datasourceVcdOrgUserRead] error retrieving org : %s", err)
	}

	user, err := adminOrg.GetUserByNameOrId(identifier, false)
	if err != nil {
		return fmt.Errorf("error retrieving user %s : %s", identifier, err)
	}

	dSet(d, "user_id", user.User.ID)
	return setOrgUserData(d, user, adminOrg)
}
