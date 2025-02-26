package v1alpha1

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/tungstenfabric/tf-operator/pkg/certificates"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _Lock sync.Mutex

const (
	SignerCAMountPath = "/etc/ssl/certs/kubernetes"
	SignerCAFilepath  = SignerCAMountPath + "/" + certificates.CAFilename
)

var signer certificates.CertificateSigner = nil

func InitCA(cl client.Client, scheme *runtime.Scheme, owner metav1.Object, ownerType string) error {
	// This might be called from reconsiles.. need sync
	_Lock.Lock()
	defer _Lock.Unlock()
	var err error
	signer, err = certificates.InitSelfCA(cl, scheme, owner, ownerType)
	if err == nil {
		return touchCertSecretsOnCAUpdate(owner.GetNamespace(), cl)
	}
	if !k8serrors.IsNotFound(err) {
		return err
	}
	if signer, err = certificates.InitK8SCA(cl, scheme, owner); err != nil {
		if !k8serrors.IsNotFound(err) {
			return err
		}
		return fmt.Errorf("Failed to init CA: neither self-signed CA secret %s nor K8S CA data %w", certificates.CaSecretName, err)
	}
	return touchCertSecretsOnCAUpdate(owner.GetNamespace(), cl)
}

// touch secrets to trigger reconcieles
func touchCertSecretsOnCAUpdate(ns string, cl client.Client) error {
	cm, err := certificates.GetCAConfigMap(ns, cl)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		return err
	}
	secretsList := &corev1.SecretList{}
	if err = cl.List(context.TODO(), secretsList, &client.ListOptions{Namespace: ns}); err != nil {
		return err
	}
	caMd5 := cm.Annotations["ca-md5"]
	for _, s := range secretsList.Items {
		if !strings.HasSuffix(s.Name, "secret-certificates") {
			continue
		}
		if s.Annotations == nil {
			s.Annotations = map[string]string{"ca-md5": ""}
		}
		if s.Annotations["ca-md5"] == caMd5 {
			continue
		}
		s.Annotations["changed-ca-md5"] = caMd5
		return cl.Update(context.TODO(), &s)
	}
	return nil
}

func retrieveDataIPs(pod corev1.Pod) []string {
	var altIPs []string
	altIP, _ := getPodDataIP(&pod)
	altIPs = append(altIPs, altIP)
	return altIPs
}

// EnsureCertificatesExist ensures pod cert is issued
func EnsureCertificatesExist(instance metav1.Object, pods []corev1.Pod, instanceType string, cl client.Client, scheme *runtime.Scheme) error {
	// This might be called from reconsiles.. need sync
	_Lock.Lock()
	defer _Lock.Unlock()
	if signer == nil {
		return fmt.Errorf("CA Signer is not initilized")
	}
	domain, err := ClusterDNSDomain(cl)
	if err != nil {
		return err
	}
	altIPs := PodAlternativeIPs{Retriever: retrieveDataIPs}
	subjects := PodsCertSubjects(domain, pods, altIPs)
	crt, err := certificates.NewCertificate(signer, cl, scheme, instance, subjects, instanceType)
	if err != nil {
		return err
	}
	return crt.EnsureExistsAndIsSigned(false)
}
