package glacier

import (
	"fmt"
	"reflect"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/aws/awsutil"
)

var (
	defaultAccountID = "-"
)

func init() {
	initRequest = func(r *aws.Request) {
		r.Handlers.Validate.PushFront(addAccountID)
		r.Handlers.Validate.PushFront(copyParams) // this happens first
		r.Handlers.Build.PushBack(addChecksum)
		r.Handlers.Build.PushBack(addAPIVersion)
	}
}

func copyParams(r *aws.Request) {
	r.Params = awsutil.CopyOf(r.Params)
}

func addAccountID(r *aws.Request) {
	if !r.ParamsFilled() {
		return
	}

	v := reflect.Indirect(reflect.ValueOf(r.Params))
	if f := v.FieldByName("AccountID"); f.Elem().String() == "" {
		f.Set(reflect.ValueOf(&defaultAccountID))
	}
}

func addChecksum(r *aws.Request) {
	if r.Body == nil {
		return
	}

	h := ComputeHashes(r.Body)

	if r.HTTPRequest.Header.Get("X-Amz-Content-Sha256") == "" {
		hstr := fmt.Sprintf("%x", h.LinearHash)
		r.HTTPRequest.Header.Set("X-Amz-Content-Sha256", hstr)
	}
	if r.HTTPRequest.Header.Get("X-Amz-Sha256-Tree-Hash") == "" {
		hstr := fmt.Sprintf("%x", h.TreeHash)
		r.HTTPRequest.Header.Set("X-Amz-Sha256-Tree-Hash", hstr)
	}
}

func addAPIVersion(r *aws.Request) {
	r.HTTPRequest.Header.Set("X-Amz-Glacier-Version", r.Service.APIVersion)
}
