package cloudflare

import (
	"context"
	"fmt"
	"log"
	"strings"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

func resourceCloudflareWorkerKV() *schema.Resource {
	return &schema.Resource{
		Schema:        resourceCloudflareWorkerKVSchema(),
		CreateContext: resourceCloudflareWorkersKVUpdate,
		ReadContext:   resourceCloudflareWorkersKVRead,
		UpdateContext: resourceCloudflareWorkersKVUpdate,
		DeleteContext: resourceCloudflareWorkersKVDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceCloudflareWorkersKVImport,
		},
	}
}

func resourceCloudflareWorkersKVRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*cloudflare.API)
	namespaceID, key := parseId(d.Id())

	value, err := client.ReadWorkersKV(ctx, namespaceID, key)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "error reading workers kv"))
	}

	if value == nil {
		d.SetId("")
		return nil
	}

	d.Set("value", string(value))
	return nil
}

func resourceCloudflareWorkersKVUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*cloudflare.API)
	namespaceID := d.Get("namespace_id").(string)
	key := d.Get("key").(string)
	value := d.Get("value").(string)

	_, err := client.WriteWorkersKV(ctx, namespaceID, key, []byte(value))
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "error creating workers kv"))
	}

	d.SetId(fmt.Sprintf("%s/%s", namespaceID, key))

	log.Printf("[INFO] Cloudflare Workers KV Namespace ID: %s", d.Id())

	return resourceCloudflareWorkersKVRead(ctx, d, meta)
}

func resourceCloudflareWorkersKVDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*cloudflare.API)
	namespaceID, key := parseId(d.Id())

	log.Printf("[INFO] Deleting Cloudflare Workers KV with id: %+v", d.Id())

	_, err := client.DeleteWorkersKV(ctx, namespaceID, key)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "error deleting workers kv"))
	}

	return nil
}

func resourceCloudflareWorkersKVImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	namespaceID, key := parseId(d.Id())

	d.Set("namespace_id", namespaceID)
	d.Set("key", key)

	resourceCloudflareWorkersKVRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}

func parseId(id string) (string, string) {
	parts := strings.SplitN(id, "/", 2)
	return parts[0], parts[1]
}
