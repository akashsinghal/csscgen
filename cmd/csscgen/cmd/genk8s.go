package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"sigs.k8s.io/yaml"
)

const (
	resolveUse = "genk8s"
)

type genk8sCmdOptions struct {
	resourceType  string
	numContainers int
	registryHost  string
	numReplicas   int
	numReferrers  int
	namespace     string
	outputPath    string
	name          string
	group         string
}

func NewCmdGenK8s(argv ...string) *cobra.Command {
	if len(argv) == 0 {
		argv = []string{os.Args[0]}
	}

	eg := fmt.Sprintf(`    # Generates a kubernetes resource template
    %s genk8s`, strings.Join(argv, " "))

	var opts genk8sCmdOptions

	cmd := &cobra.Command{
		Use:     resolveUse,
		Short:   "Generates a kubernetes resource template",
		Example: eg,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return createResource(opts)
		},
	}

	flags := cmd.Flags()

	flags.StringVarP(&opts.resourceType, "resource-type", "t", "deployment", "Resource Type (deployment, job, pod)")
	flags.StringVar(&opts.registryHost, "registry-host", "docker.io", "Registry Host")
	flags.IntVarP(&opts.numContainers, "num-containers", "c", 1, "Number of containers")
	flags.IntVar(&opts.numReplicas, "num-replicas", 1, "Number of replicas")
	flags.IntVar(&opts.numReferrers, "num-referrers", 1, "Number of referrers")
	flags.StringVarP(&opts.namespace, "namespace", "n", "", "Namespace")
	flags.StringVarP(&opts.outputPath, "output-file", "f", "", "Output file name")
	flags.StringVar(&opts.name, "name", "", "Name")
	flags.StringVar(&opts.group, "group", "", "Group")
	return cmd
}

func createResource(opts genk8sCmdOptions) error {
	imageName := getImageName(opts.numContainers, opts.numReferrers)
	returnTemplate := ""
	var err error
	castedNumReplicas := int32(opts.numReplicas)
	switch opts.resourceType {
	case "deployment":
		name := opts.name
		group := opts.group
		if name == "" {
			name = "{{.Name}}"
		}
		if group == "" {
			group = "{{.Group}}"
		}
		returnTemplate, err = createDeployment(&castedNumReplicas, opts.numContainers, imageName, opts.registryHost, opts.namespace, name, group)
		if err != nil {
			return err
		}
	case "job":
		name := opts.name
		group := opts.group
		if name == "" {
			name = "{{.Name}}"
		}
		if group == "" {
			group = "{{.Group}}"
		}
		returnTemplate, err = createJob(&castedNumReplicas, opts.namespace, name, opts.numContainers, imageName, opts.registryHost, group)
		if err != nil {
			return err
		}
	// case "pod":
	// 	return createPod(numReplicas, numContainers, imageName, registryName, namespace)
	default:
		return fmt.Errorf("invalid resource type: %s", opts.resourceType)
	}

	if opts.outputPath != "" {
		f, err := os.Create(opts.outputPath)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = f.WriteString(returnTemplate)
		if err != nil {
			return err
		}
	} else {
		fmt.Println(returnTemplate)
	}
	return nil
}

func createDeployment(numReplicas *int32, numContainers int, imageName string, registryName string, namespace string, name string, group string) (string, error) {
	containers := make([]corev1.Container, numContainers)
	for i := 0; i < numContainers; i++ {
		containers[i] = corev1.Container{
			Name:  fmt.Sprintf("%s%v", imageName, i+1),
			Image: fmt.Sprintf("%s/%s:%v", registryName, imageName, i+1),
		}
	}
	var objectMeta metav1.ObjectMeta
	if namespace != "" {
		objectMeta = metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"group": group,
			},
		}
	} else {
		objectMeta = metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"group": group,
			},
		}
	}
	template := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: objectMeta,
		Spec: appsv1.DeploymentSpec{
			Replicas: numReplicas,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"name":  name,
						"group": group,
					},
				},
				Spec: corev1.PodSpec{
					Containers: containers,
				},
			},
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"name": name,
				},
			},
		},
	}

	bytes, err := yaml.Marshal(template)
	if err != nil {
		return "", fmt.Errorf("failed to create deployment: %v", err)
	}
	return string(bytes), nil
}

func createJob(numReplicas *int32, namespace string, jobName string, numContainers int, imageName string, registryName string, group string) (string, error) {
	containers := make([]corev1.Container, numContainers)
	for i := 0; i < numContainers; i++ {
		containers[i] = corev1.Container{
			Name:  fmt.Sprintf("%s%v", imageName, i+1),
			Image: fmt.Sprintf("%s/%s:%v", registryName, imageName, i+1),
		}
	}
	var objectMeta metav1.ObjectMeta
	if namespace != "" {
		objectMeta = metav1.ObjectMeta{
			Name:      jobName,
			Namespace: namespace,
			Labels: map[string]string{
				"group": group,
			},
		}
	} else {
		objectMeta = metav1.ObjectMeta{
			Name: jobName,
			Labels: map[string]string{
				"group": group,
			},
		}
	}
	jobCompletionMode := batchv1.IndexedCompletion
	template := &batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: "batch/v1",
		},
		ObjectMeta: objectMeta,
		Spec: batchv1.JobSpec{
			Parallelism:    numReplicas,
			Completions:    numReplicas,
			CompletionMode: &jobCompletionMode,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"group": group,
					},
				},
				Spec: corev1.PodSpec{
					Containers:    containers,
					RestartPolicy: corev1.RestartPolicyOnFailure,
				},
			},
		},
	}
	bytes, err := yaml.Marshal(template)
	if err != nil {
		return "", fmt.Errorf("failed to create job: %v", err)
	}
	return string(bytes), nil
}

func getImageName(numContainers int, numReferrers int) string {
	return fmt.Sprintf("%d-containers-%d-referrers", numContainers, numReferrers)
}
