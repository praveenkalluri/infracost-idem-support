package aws

import (
	"github.com/infracost/infracost/internal/resources/aws"

	"github.com/infracost/infracost/internal/schema"
)

func GetS3BucketRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name: "states.aws.s3.bucket.present",
		Notes: []string{
			"S3 replication time control data transfer, and batch operations are not supported by Terraform.",
		},
		CoreRFunc: NewS3BucketResource,
		ReferenceAttributes: []string{
			"states.aws.s3.bucket_lifecycle.present:bucket",
			"states.aws.cloudfront.distribution.present:origins.0.DomainName",
			"states.aws.cloudfront.distribution.present:origins.0.Id",
		},
	}
}

func NewS3BucketResource(d *schema.ResourceData) schema.CoreResource {
	storageClassNames := map[string]string{
		"STANDARD":            "standard",
		"INTELLIGENT_TIERING": "intelligent_tiering",
		"STANDARD_IA":         "standard_infrequent_access",
		"ONEZONE_IA":          "one_zone_infrequent_access",
		"GLACIER":             "glacier_flexible_retrieval",
		"DEEP_ARCHIVE":        "glacier_deep_archive",
	}

	objTagsEnabled := false

	// Always add the standard storage class
	lifecycleStorageClassMap := map[string]bool{
		"standard": true,
	}

	// TODO: As Idem do not support adding life cycle rules in the same bucket state
	// TODO: we need to refer it from another class and Implement this logic
	for _, rule := range d.Get("lifecycle_rule").Array() {
		if !rule.Get("enabled").Bool() {
			continue
		}

		if len(rule.Get("tags").Map()) > 0 {
			objTagsEnabled = true
		}

		for _, t := range rule.Get("transition").Array() {
			storageClass := storageClassNames[t.Get("storage_class").String()]
			if _, ok := lifecycleStorageClassMap[storageClass]; !ok && storageClass != "" {
				lifecycleStorageClassMap[storageClass] = true
			}
		}

		for _, t := range rule.Get("noncurrent_version_transition").Array() {
			if !rule.Get("enabled").Bool() {
				continue
			}
			storageClass := storageClassNames[t.Get("storage_class").String()]
			if _, ok := lifecycleStorageClassMap[storageClass]; !ok && storageClass != "" {
				lifecycleStorageClassMap[storageClass] = true
			}
		}
	}

	lifecycleStorageClasses := make([]string, 0, len(lifecycleStorageClassMap))
	for storageClass := range lifecycleStorageClassMap {
		lifecycleStorageClasses = append(lifecycleStorageClasses, storageClass)
	}

	return &aws.S3Bucket{
		Address:                 d.Address,
		Region:                  d.Get("region").String(),
		Name:                    d.Get("bucket").String(),
		ObjectTagsEnabled:       objTagsEnabled,
		LifecycleStorageClasses: lifecycleStorageClasses,
	}
}
