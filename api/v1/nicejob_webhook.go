/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	"github.com/robfig/cron"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	validationutils "k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var nicejoblog = logf.Log.WithName("nicejob-resource")

func (r *NiceJob) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-webapp-tutorial-io-v1-nicejob,mutating=true,failurePolicy=fail,sideEffects=None,groups=webapp.tutorial.io,resources=nicejobs,verbs=create;update,versions=v1,name=mnicejob.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Defaulter = &NiceJob{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *NiceJob) Default() {
	nicejoblog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
	if r.Spec.ConcurrencyPolicy == "" {
		r.Spec.ConcurrencyPolicy = AllowConcurrent
	}
	if r.Spec.Suspend == nil {
		r.Spec.Suspend = new(bool)
	}
	if r.Spec.SuccessfulJobsHistoryLimit == nil {
		r.Spec.SuccessfulJobsHistoryLimit = new(int32)
		*r.Spec.SuccessfulJobsHistoryLimit = 3
	}
	if r.Spec.FailedJobsHistoryLimit == nil {
		r.Spec.FailedJobsHistoryLimit = new(int32)
		*r.Spec.FailedJobsHistoryLimit = 1
	}
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-webapp-tutorial-io-v1-nicejob,mutating=false,failurePolicy=fail,sideEffects=None,groups=webapp.tutorial.io,resources=nicejobs,verbs=create;update,versions=v1,name=vnicejob.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &NiceJob{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *NiceJob) ValidateCreate() error {
	nicejoblog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return r.validateNiceJob()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *NiceJob) ValidateUpdate(old runtime.Object) error {
	nicejoblog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return r.validateNiceJob()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *NiceJob) ValidateDelete() error {
	nicejoblog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}

/*
We validate the name and the spec of the CronJob.
*/

func (r *NiceJob) validateNiceJob() error {
	var allErrs field.ErrorList
	if err := r.validateNiceJobName(); err != nil {
		allErrs = append(allErrs, err)
	}
	if err := r.validateNiceJobSpec(); err != nil {
		allErrs = append(allErrs, err)
	}
	if len(allErrs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		schema.GroupKind{Group: "batch.tutorial.kubebuilder.io", Kind: "CronJob"},
		r.Name, allErrs)
}

/*
Some fields are declaratively validated by OpenAPI schema.
You can find kubebuilder validation markers (prefixed
with `// +kubebuilder:validation`) in the
[Designing an API](api-design.md) section.
You can find all of the kubebuilder supported markers for
declaring validation by running `controller-gen crd -w`,
or [here](/reference/markers/crd-validation.md).
*/

func (r *NiceJob) validateNiceJobSpec() *field.Error {
	// The field helpers from the kubernetes API machinery help us return nicely
	// structured validation errors.
	return validateScheduleFormat(
		r.Spec.Schedule,
		field.NewPath("spec").Child("schedule"))
}

/*
We'll need to validate the [cron](https://en.wikipedia.org/wiki/Cron) schedule
is well-formatted.
*/

func validateScheduleFormat(schedule string, fldPath *field.Path) *field.Error {
	if _, err := cron.ParseStandard(schedule); err != nil {
		return field.Invalid(fldPath, schedule, err.Error())
	}
	return nil
}

/*
Validating the length of a string field can be done declaratively by
the validation schema.
But the `ObjectMeta.Name` field is defined in a shared package under
the apimachinery repo, so we can't declaratively validate it using
the validation schema.
*/

func (r *NiceJob) validateNiceJobName() *field.Error {
	if len(r.ObjectMeta.Name) > validationutils.DNS1035LabelMaxLength-11 {
		// The job name length is 63 character like all Kubernetes objects
		// (which must fit in a DNS subdomain). The cronjob controller appends
		// a 11-character suffix to the cronjob (`-$TIMESTAMP`) when creating
		// a job. The job name length limit is 63 characters. Therefore cronjob
		// names must have length <= 63-11=52. If we don't validate this here,
		// then job creation will fail later.
		return field.Invalid(field.NewPath("metadata").Child("name"), r.Name, "must be no more than 52 characters")
	}
	return nil
}

// +kubebuilder:docs-gen:collapse=Validate object name
