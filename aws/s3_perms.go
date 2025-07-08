package aws

// ACLTypes represents the different types of Access Control Lists (ACLs) for objects in a system.
//
// Usage example:
//
// // Custom String() method for ACLTypes
//
//	func (a *ACLTypes) String() string {
//	    return string(*a)
//	}
//
// // Define ACLTypes constants
// S3ObjectCannedAclPublicRead ACLTypes = "public-read"
// S3ObjectCannedAclAuthenticatedRead ACLTypes = "authenticated-read"
// S3ObjectCannedAclBucketOwnerRead ACLTypes = "bucket-owner-read"
// S3ObjectCannedAclBucketOwnerFullControl ACLTypes = "bucket-owner-full-control"
//
// // Struct using ACLTypes
//
//	WithACL struct {
//	    ACLName            ACLTypes
//	    IsAclTypeOperation bool
//	}
type ACLTypes string

// String returns the string representation of the ACLTypes value.
func (a *ACLTypes) String() string {
	return string(*a)
}

const (
	// S3ObjectCannedAclPublicRead is a S3ObjectCannedAcl enum value
	S3ObjectCannedAclPublicRead ACLTypes = "public-read"

	// S3ObjectCannedAclAuthenticatedRead is a S3ObjectCannedAcl enum value
	S3ObjectCannedAclAuthenticatedRead ACLTypes = "authenticated-read"

	// S3ObjectCannedAclBucketOwnerRead is a S3ObjectCannedAcl enum value
	S3ObjectCannedAclBucketOwnerRead ACLTypes = "bucket-owner-read"

	// S3ObjectCannedAclBucketOwnerFullControl is a S3ObjectCannedAcl enum value
	S3ObjectCannedAclBucketOwnerFullControl ACLTypes = "bucket-owner-full-control"
)

// WithACL represents a struct that holds information about Access Control Lists (ACLs) for objects in a system.
// ACLName represents the ACL type and is of type ACLTypes.
// IsAclTypeOperation represents a boolean value indicating whether ACL operation has been performed.
type WithACL struct {
	ACLName            ACLTypes
	IsAclTypeOperation bool
}

// WithMetadata represents an object that has associated metadata.
// It contains a `Metadata` map that stores key-value pairs, where the key is a string and the value is a pointer to a string.
// The `IsMetadata` field indicates whether the object's metadata should be included in the operation.
//
// Usage example:
// s3 := &S3{}
// s3.MetadataOperation = true
// s3.Metadata = make(map[string]*string)
// s3.Metadata["key1"] = aws.String("value1")
// s3.Metadata["key2"] = aws.String("value2")
//
//	wm := &WithMetadata{
//	   Metadata:   s3.Metadata,
//	   IsMetadata: s3.MetadataOperation,
//	}
//
// s3.PutObjectInBucket("bucketName", "key", []byte("object"), wm)
type WithMetadata struct {
	Metadata   map[string]*string
	IsMetadata bool
}

// WithContentType represents an object that has a content type.
// Usage example:
//
//	func (w *WithContentType) GetS3PermBool(s3 *S3) {
//	    s3.ContentType = w.ContentType
//	}
//
// WithContentType struct contains the ContentType field.
type WithContentType struct {
	ContentType string
}

// GetS3PermBool sets the MetadataOperation and Metadata fields of the S3 object
// passed as an argument to WithMetadata's IsMetadata and Metadata fields.
func (w *WithMetadata) GetS3PermBool(s3 *S3) {
	s3.MetadataOperation = w.IsMetadata
	s3.Metadata = w.Metadata
}

// GetS3PermBool sets the ContentType for the given S3 object based on the ContentType specified in WithContentType.
func (w *WithContentType) GetS3PermBool(s3 *S3) {
	s3.ContentType = w.ContentType
}

// GetS3PermBool sets the ACLOperation and ACLName fields of the S3 struct based on the values of the IsAclTypeOperation and ACLName fields of the WithACL struct.
// The S3 struct is passed as a parameter and modified in place.
// If IsAclTypeOperation is true, ACLOperation is set to true.
// ACLName is set to the string representation of ACLName using its String method.
// Please refer to the WithACL and S3 structs for more information on the fields used in this function.
func (w *WithACL) GetS3PermBool(s3 *S3) {
	s3.ACLOperation = w.IsAclTypeOperation
	s3.ACLName = w.ACLName.String()
}
