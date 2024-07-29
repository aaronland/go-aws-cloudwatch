// Code generated by smithy-go-codegen DO NOT EDIT.

package iam

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Changes the status of the specified user signing certificate from active to
// disabled, or vice versa. This operation can be used to disable an IAM user's
// signing certificate as part of a certificate rotation work flow.
//
// If the UserName field is not specified, the user name is determined implicitly
// based on the Amazon Web Services access key ID used to sign the request. This
// operation works for access keys under the Amazon Web Services account.
// Consequently, you can use this operation to manage Amazon Web Services account
// root user credentials even if the Amazon Web Services account has no associated
// users.
func (c *Client) UpdateSigningCertificate(ctx context.Context, params *UpdateSigningCertificateInput, optFns ...func(*Options)) (*UpdateSigningCertificateOutput, error) {
	if params == nil {
		params = &UpdateSigningCertificateInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "UpdateSigningCertificate", params, optFns, c.addOperationUpdateSigningCertificateMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*UpdateSigningCertificateOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type UpdateSigningCertificateInput struct {

	// The ID of the signing certificate you want to update.
	//
	// This parameter allows (through its [regex pattern]) a string of characters that can consist of
	// any upper or lowercased letter or digit.
	//
	// [regex pattern]: http://wikipedia.org/wiki/regex
	//
	// This member is required.
	CertificateId *string

	//  The status you want to assign to the certificate. Active means that the
	// certificate can be used for programmatic calls to Amazon Web Services Inactive
	// means that the certificate cannot be used.
	//
	// This member is required.
	Status types.StatusType

	// The name of the IAM user the signing certificate belongs to.
	//
	// This parameter allows (through its [regex pattern]) a string of characters consisting of upper
	// and lowercase alphanumeric characters with no spaces. You can also include any
	// of the following characters: _+=,.@-
	//
	// [regex pattern]: http://wikipedia.org/wiki/regex
	UserName *string

	noSmithyDocumentSerde
}

type UpdateSigningCertificateOutput struct {
	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationUpdateSigningCertificateMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsAwsquery_serializeOpUpdateSigningCertificate{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsAwsquery_deserializeOpUpdateSigningCertificate{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "UpdateSigningCertificate"); err != nil {
		return fmt.Errorf("add protocol finalizers: %v", err)
	}

	if err = addlegacyEndpointContextSetter(stack, options); err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = addClientRequestID(stack); err != nil {
		return err
	}
	if err = addComputeContentLength(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = addComputePayloadSHA256(stack); err != nil {
		return err
	}
	if err = addRetry(stack, options); err != nil {
		return err
	}
	if err = addRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = addRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addClientUserAgent(stack, options); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = addSetLegacyContextSigningOptionsMiddleware(stack); err != nil {
		return err
	}
	if err = addTimeOffsetBuild(stack, c); err != nil {
		return err
	}
	if err = addUserAgentRetryMode(stack, options); err != nil {
		return err
	}
	if err = addOpUpdateSigningCertificateValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opUpdateSigningCertificate(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = addRecursionDetection(stack); err != nil {
		return err
	}
	if err = addRequestIDRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addResponseErrorMiddleware(stack); err != nil {
		return err
	}
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	if err = addDisableHTTPSMiddleware(stack, options); err != nil {
		return err
	}
	return nil
}

func newServiceMetadataMiddleware_opUpdateSigningCertificate(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "UpdateSigningCertificate",
	}
}
