package internal

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func NewK8sClient() (string, *kubernetes.Clientset) {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)

	if err != nil {
		fmt.Printf("Can't build config from kubeconfig \n")
		fmt.Printf("%s \n", err.Error())
		os.Exit(2)
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Can't login kubernetes cluster \n")
		fmt.Printf("%s \n", err.Error())
		os.Exit(2)
	}

	return config.Host, clientset
}

func GetTokenServiceAccount(serviceAccount string, namespace string, clientset *kubernetes.Clientset) string {

	sn, err := clientset.CoreV1().ServiceAccounts(namespace).Get(context.TODO(), serviceAccount, metav1.GetOptions{})
	if err != nil {
		fmt.Printf("Can't get service account: %s \n", serviceAccount)
		fmt.Printf("%s \n", err.Error())
		os.Exit(2)
	}
	tn := sn.Secrets[0].DeepCopy().Name

	s, err := clientset.CoreV1().Secrets(namespace).Get(context.TODO(), tn, metav1.GetOptions{})
	if err != nil {
		fmt.Printf("Can't get secret from service account %s \n", serviceAccount)
		fmt.Printf("%s \n", err.Error())
		os.Exit(2)
	}
	t := string(s.DeepCopy().Data["token"])

	return t
}
